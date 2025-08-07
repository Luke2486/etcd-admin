package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"etcd-admin-backend/internal/models"
	"etcd-admin-backend/internal/services"
	"etcd-admin-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

// KVHandler KV操作处理器
type KVHandler struct {
	etcdService *services.EtcdService
}

// NewKVHandler 创建KV处理器
func NewKVHandler(etcdService *services.EtcdService) *KVHandler {
	return &KVHandler{
		etcdService: etcdService,
	}
}

// SetValueRequest 设置值请求
type SetValueRequest struct {
	Value interface{} `json:"value" binding:"required"`
}

// ListKeysResponse 列出键响应
type ListKeysResponse struct {
	Keys []string `json:"keys"`
}

// GetValueResponse 获取值响应
type GetValueResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ListKeys 列出所有键
func (h *KVHandler) ListKeys(c *gin.Context) {
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

	// 从etcd获取键列表
	keys, err := h.etcdService.ListKeys(&connection, prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to list keys",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Keys retrieved successfully",
		"data": ListKeysResponse{
			Keys: keys,
		},
	})
}

// GetValue 获取键值
func (h *KVHandler) GetValue(c *gin.Context) {
	connectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection_id",
		})
		return
	}

	key := c.Param("key")
	// 移除通配符參數前的斜杠，因為我們使用 *key 路由
	if key != "" && key[0] == '/' {
		key = key[1:]
	}

	// 重新添加開頭的斜杠，因為etcd的key通常以斜杠開頭
	if key != "" && key[0] != '/' {
		key = "/" + key
	}

	if key == "" || key == "/" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Key is required",
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

	// 从etcd获取值
	value, err := h.etcdService.GetValue(&connection, key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Key not found or failed to get value",
			"error":   err.Error(),
		})
		return
	}

	// 尝试解析JSON
	var jsonValue interface{}
	if err := json.Unmarshal([]byte(value), &jsonValue); err != nil {
		// 如果不是JSON，返回原始字符串
		jsonValue = value
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Value retrieved successfully",
		"data": GetValueResponse{
			Key:   key,
			Value: jsonValue,
		},
	})
}

// SetValue 设置键值
func (h *KVHandler) SetValue(c *gin.Context) {
	connectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection_id",
		})
		return
	}

	key := c.Param("key")
	// 移除通配符參數前的斜杠，因為我們使用 *key 路由
	if key != "" && key[0] == '/' {
		key = key[1:]
	}

	// 重新添加開頭的斜杠，因為etcd的key通常以斜杠開頭
	if key != "" && key[0] != '/' {
		key = "/" + key
	}

	if key == "" || key == "/" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Key is required",
		})
		return
	}

	var req SetValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
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

	// 检查连接是否为只读
	if connection.IsReadOnly {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Connection is read-only, cannot set values",
		})
		return
	}

	// 将值转换为JSON字符串
	valueBytes, err := json.Marshal(req.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to serialize value",
		})
		return
	}

	// 设置到etcd
	if err := h.etcdService.SetValue(&connection, key, string(valueBytes)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to set value",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Value set successfully",
		"data": map[string]interface{}{
			"key":   key,
			"value": req.Value,
		},
	})
}

// DeleteKey 删除键
func (h *KVHandler) DeleteKey(c *gin.Context) {
	connectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection_id",
		})
		return
	}

	key := c.Param("key")
	// 移除通配符參數前的斜杠，因為我們使用 *key 路由
	if key != "" && key[0] == '/' {
		key = key[1:]
	}

	// 重新添加開頭的斜杠，因為etcd的key通常以斜杠開頭
	if key != "" && key[0] != '/' {
		key = "/" + key
	}

	if key == "" || key == "/" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Key is required",
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

	// 检查连接是否为只读
	if connection.IsReadOnly {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "Connection is read-only, cannot delete keys",
		})
		return
	}

	// 从etcd删除键
	if err := h.etcdService.DeleteKey(&connection, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete key",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Key deleted successfully",
		"data": map[string]string{
			"key": key,
		},
	})
}
