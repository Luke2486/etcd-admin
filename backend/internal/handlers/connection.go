package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"etcd-admin-backend/internal/models"
	"etcd-admin-backend/internal/services"
	"etcd-admin-backend/pkg/database"
)

// ConnectionHandler 连接管理处理器
type ConnectionHandler struct {
	etcdService *services.EtcdService
}

// NewConnectionHandler 创建连接处理器
func NewConnectionHandler(etcdService *services.EtcdService) *ConnectionHandler {
	return &ConnectionHandler{
		etcdService: etcdService,
	}
}

// CreateConnectionRequest 创建连接请求
type CreateConnectionRequest struct {
	Name        string   `json:"name" binding:"required"`
	Endpoints   []string `json:"endpoints" binding:"required"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Description string   `json:"description"`
	IsReadOnly  bool     `json:"is_readonly"`
}

// UpdateConnectionRequest 更新连接请求
type UpdateConnectionRequest struct {
	Name        string   `json:"name" binding:"required"`
	Endpoints   []string `json:"endpoints" binding:"required"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Description string   `json:"description"`
	IsReadOnly  bool     `json:"is_readonly"`
}

// ListConnections 获取连接列表
func (h *ConnectionHandler) ListConnections(c *gin.Context) {
	var connections []models.Connection

	if err := database.GetDB().Find(&connections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch connections",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Connections retrieved successfully",
		"data":    connections,
	})
}

// CreateConnection 创建新连接
func (h *ConnectionHandler) CreateConnection(c *gin.Context) {
	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 将endpoints数组转换为JSON字符串
	endpointsJSON, err := json.Marshal(req.Endpoints)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid endpoints format",
		})
		return
	}

	connection := models.Connection{
		Name:        req.Name,
		Endpoints:   string(endpointsJSON),
		Username:    req.Username,
		Password:    req.Password,
		Description: req.Description,
		IsActive:    true,
		IsReadOnly:  req.IsReadOnly,
	}

	// 测试连接
	if err := h.etcdService.TestConnection(&connection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to connect to etcd",
			"error":   err.Error(),
		})
		return
	}

	if err := database.GetDB().Create(&connection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create connection",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Connection created successfully",
		"data":    connection,
	})
}

// GetConnection 获取单个连接
func (h *ConnectionHandler) GetConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection ID",
		})
		return
	}

	var connection models.Connection
	if err := database.GetDB().First(&connection, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Connection not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Connection retrieved successfully",
		"data":    connection,
	})
}

// UpdateConnection 更新连接
func (h *ConnectionHandler) UpdateConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection ID",
		})
		return
	}

	var req UpdateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 查找现有连接
	var connection models.Connection
	if err := database.GetDB().First(&connection, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Connection not found",
		})
		return
	}

	// 将endpoints数组转换为JSON字符串
	endpointsJSON, err := json.Marshal(req.Endpoints)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid endpoints format",
		})
		return
	}

	// 创建临时连接用于测试
	testConnection := models.Connection{
		Name:        req.Name,
		Endpoints:   string(endpointsJSON),
		Username:    req.Username,
		Password:    req.Password,
		Description: req.Description,
		IsReadOnly:  req.IsReadOnly,
		TLSEnabled:  connection.TLSEnabled, // 保持原有TLS设置
		CertFile:    connection.CertFile,
		KeyFile:     connection.KeyFile,
		CAFile:      connection.CAFile,
	}

	// 测试连接
	if err := h.etcdService.TestConnection(&testConnection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to connect to etcd with new settings",
			"error":   err.Error(),
		})
		return
	}

	// 关闭旧的etcd客户端连接
	h.etcdService.CloseClient(uint(id))

	// 更新连接信息
	connection.Name = req.Name
	connection.Endpoints = string(endpointsJSON)
	connection.Username = req.Username
	connection.Password = req.Password
	connection.Description = req.Description
	connection.IsReadOnly = req.IsReadOnly

	if err := database.GetDB().Save(&connection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update connection",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Connection updated successfully",
		"data":    connection,
	})
}

// DeleteConnection 删除连接
func (h *ConnectionHandler) DeleteConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection ID",
		})
		return
	}

	// 关闭客户端连接
	h.etcdService.CloseClient(uint(id))

	if err := database.GetDB().Delete(&models.Connection{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete connection",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Connection deleted successfully",
	})
}

// TestConnection 测试连接
func (h *ConnectionHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid connection ID",
		})
		return
	}

	var connection models.Connection
	if err := database.GetDB().First(&connection, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Connection not found",
		})
		return
	}

	if err := h.etcdService.TestConnection(&connection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Connection test failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Connection test successful",
	})
}
