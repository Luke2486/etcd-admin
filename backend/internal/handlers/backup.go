package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"etcd-admin-backend/internal/models"
	"etcd-admin-backend/internal/services"
	"etcd-admin-backend/pkg/database"
)

// BackupHandler 备份处理器
type BackupHandler struct {
	etcdService *services.EtcdService
}

// NewBackupHandler 创建备份处理器
func NewBackupHandler(etcdService *services.EtcdService) *BackupHandler {
	return &BackupHandler{
		etcdService: etcdService,
	}
}

// BackupData 备份数据结构
type BackupData struct {
	ConnectionName string                 `json:"connection_name"`
	ConnectionID   uint                   `json:"connection_id"`
	ExportTime     time.Time              `json:"export_time"`
	Data           map[string]interface{} `json:"data"`
}

// ImportRequest 导入请求
type ImportRequest struct {
	Data      map[string]interface{} `json:"data" binding:"required"`
	Overwrite bool                   `json:"overwrite"` // 是否覆盖已存在的key
}

// ExportBackup 导出备份
func (h *BackupHandler) ExportBackup(c *gin.Context) {
	connectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection_id",
		})
		return
	}

	// 获取连接配置
	var connection models.Connection
	if err := database.GetDB().First(&connection, uint(connectionID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Connection not found",
		})
		return
	}

	// 获取前缀参数
	prefix := c.DefaultQuery("prefix", "")

	// 从etcd获取所有KV数据
	kvData, err := h.etcdService.GetAllKV(&connection, prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to export data",
			"error":   err.Error(),
		})
		return
	}

	// 转换为适合导出的格式
	exportData := make(map[string]interface{})
	for key, value := range kvData {
		// 尝试解析JSON
		var jsonValue interface{}
		if err := json.Unmarshal([]byte(value), &jsonValue); err != nil {
			// 如果不是JSON，存储原始字符串
			exportData[key] = value
		} else {
			exportData[key] = jsonValue
		}
	}

	backup := BackupData{
		ConnectionName: connection.Name,
		ConnectionID:   connection.ID,
		ExportTime:     time.Now(),
		Data:           exportData,
	}

	// 设置下载文件名
	filename := fmt.Sprintf("etcd-backup-%s-%s.json",
		connection.Name,
		time.Now().Format("20060102-150405"))

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, backup)
}

// ImportBackup 导入备份
func (h *BackupHandler) ImportBackup(c *gin.Context) {
	connectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection_id",
		})
		return
	}

	// 获取连接配置
	var connection models.Connection
	if err := database.GetDB().First(&connection, uint(connectionID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Connection not found",
		})
		return
	}

	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	successCount := 0
	errorCount := 0
	errors := make([]string, 0)

	// 导入数据到etcd
	for key, value := range req.Data {
		// 如果不覆盖，先检查key是否存在
		if !req.Overwrite {
			if _, err := h.etcdService.GetValue(&connection, key); err == nil {
				// key已存在，跳过
				continue
			}
		}

		// 将值转换为JSON字符串
		valueBytes, err := json.Marshal(value)
		if err != nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("Failed to serialize value for key %s: %v", key, err))
			continue
		}

		// 设置到etcd
		if err := h.etcdService.SetValue(&connection, key, string(valueBytes)); err != nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("Failed to set key %s: %v", key, err))
			continue
		}

		successCount++
	}

	response := gin.H{
		"status":        "success",
		"message":       "Import completed",
		"success_count": successCount,
		"error_count":   errorCount,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}
