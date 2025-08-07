package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"etcd-admin-backend/internal/models"
	"etcd-admin-backend/internal/services"
	"etcd-admin-backend/pkg/database"
)

// TransferHandler KV传输处理器
type TransferHandler struct {
	etcdService *services.EtcdService
}

// NewTransferHandler 创建传输处理器
func NewTransferHandler(etcdService *services.EtcdService) *TransferHandler {
	return &TransferHandler{
		etcdService: etcdService,
	}
}

// TransferRequest 传输请求
type TransferRequest struct {
	SourceConnectionID uint     `json:"source_connection_id" binding:"required"`
	TargetConnectionID uint     `json:"target_connection_id" binding:"required"`
	Keys               []string `json:"keys"`          // 指定要传输的keys，为空则传输所有
	Prefix             string   `json:"prefix"`        // 前缀过滤
	Overwrite          bool     `json:"overwrite"`     // 是否覆盖目标中已存在的key
	KeyMapping         bool     `json:"key_mapping"`   // 是否启用键名映射
	SourcePrefix       string   `json:"source_prefix"` // 源前缀
	TargetPrefix       string   `json:"target_prefix"` // 目标前缀
}

// TransferResponse 传输响应
type TransferResponse struct {
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	SkippedCount int      `json:"skipped_count"`
	Errors       []string `json:"errors,omitempty"`
	Details      []string `json:"details,omitempty"`
}

// TransferKV 在连接间传输KV数据
func (h *TransferHandler) TransferKV(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 验证源连接和目标连接不能相同
	if req.SourceConnectionID == req.TargetConnectionID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Source and target connections cannot be the same",
		})
		return
	}

	// 获取源连接配置
	var sourceConnection models.Connection
	if err := database.GetDB().First(&sourceConnection, req.SourceConnectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Source connection not found",
		})
		return
	}

	// 获取目标连接配置
	var targetConnection models.Connection
	if err := database.GetDB().First(&targetConnection, req.TargetConnectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Target connection not found",
		})
		return
	}

	response := TransferResponse{
		Errors:  make([]string, 0),
		Details: make([]string, 0),
	}

	var keysToTransfer []string
	var err error

	// 确定要传输的键列表
	if len(req.Keys) > 0 {
		// 传输指定的键
		keysToTransfer = req.Keys
	} else {
		// 传输所有键或按前缀过滤的键
		prefix := req.Prefix
		if req.KeyMapping && req.SourcePrefix != "" {
			prefix = req.SourcePrefix
		}

		keysToTransfer, err = h.etcdService.ListKeys(&sourceConnection, prefix)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to list keys from source connection",
				"error":   err.Error(),
			})
			return
		}
	}

	// 传输每个键
	for _, key := range keysToTransfer {
		// 从源连接获取值
		value, err := h.etcdService.GetValue(&sourceConnection, key)
		if err != nil {
			response.ErrorCount++
			response.Errors = append(response.Errors,
				"Failed to get key '"+key+"' from source: "+err.Error())
			continue
		}

		// 确定目标键名
		targetKey := key
		if req.KeyMapping && req.SourcePrefix != "" && req.TargetPrefix != "" {
			// 进行键名映射：将源前缀替换为目标前缀
			if len(key) >= len(req.SourcePrefix) && key[:len(req.SourcePrefix)] == req.SourcePrefix {
				targetKey = req.TargetPrefix + key[len(req.SourcePrefix):]
			}
		}

		// 如果不覆盖，检查目标键是否已存在
		if !req.Overwrite {
			if _, err := h.etcdService.GetValue(&targetConnection, targetKey); err == nil {
				response.SkippedCount++
				response.Details = append(response.Details,
					"Skipped existing key: "+targetKey)
				continue
			}
		}

		// 设置到目标连接
		if err := h.etcdService.SetValue(&targetConnection, targetKey, value); err != nil {
			response.ErrorCount++
			response.Errors = append(response.Errors,
				"Failed to set key '"+targetKey+"' to target: "+err.Error())
			continue
		}

		response.SuccessCount++
		response.Details = append(response.Details,
			"Successfully transferred: "+key+" -> "+targetKey)
	}

	status := "success"
	message := "Transfer completed successfully"
	if response.ErrorCount > 0 {
		status = "partial_success"
		message = "Transfer completed with some errors"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": message,
		"data":    response,
	})
}

// CopyKey 复制单个键
func (h *TransferHandler) CopyKey(c *gin.Context) {
	sourceConnID, err := strconv.ParseUint(c.Query("source_connection_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid source_connection_id",
		})
		return
	}

	targetConnID, err := strconv.ParseUint(c.Query("target_connection_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid target_connection_id",
		})
		return
	}

	sourceKey := c.Param("key")
	if sourceKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Key is required",
		})
		return
	}

	targetKey := c.DefaultQuery("target_key", sourceKey)
	overwrite := c.DefaultQuery("overwrite", "false") == "true"

	// 获取连接配置
	var sourceConnection, targetConnection models.Connection
	if err := database.GetDB().First(&sourceConnection, uint(sourceConnID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Source connection not found",
		})
		return
	}

	if err := database.GetDB().First(&targetConnection, uint(targetConnID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Target connection not found",
		})
		return
	}

	// 从源获取值
	value, err := h.etcdService.GetValue(&sourceConnection, sourceKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Failed to get value from source",
			"error":   err.Error(),
		})
		return
	}

	// 检查目标是否已存在
	if !overwrite {
		if _, err := h.etcdService.GetValue(&targetConnection, targetKey); err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Target key already exists",
			})
			return
		}
	}

	// 设置到目标
	if err := h.etcdService.SetValue(&targetConnection, targetKey, value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to set value to target",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Key copied successfully",
		"data": map[string]interface{}{
			"source_key": sourceKey,
			"target_key": targetKey,
			"value":      value,
		},
	})
}
