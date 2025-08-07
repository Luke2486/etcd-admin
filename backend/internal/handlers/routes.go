package handlers

import (
	"github.com/gin-gonic/gin"

	"etcd-admin-backend/internal/config"
	"etcd-admin-backend/internal/middleware"
	"etcd-admin-backend/internal/services"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, cfg *config.Config, etcdService *services.EtcdService) {
	// 创建处理器
	authHandler := NewAuthHandler(cfg)
	connectionHandler := NewConnectionHandler(etcdService)
	kvHandler := NewKVHandler(etcdService)
	backupHandler := NewBackupHandler(etcdService)
	transferHandler := NewTransferHandler(etcdService)

	// API路由组
	api := r.Group("/api/v1")

	// 公开路由（不需要认证）
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)
	}

	// 需要认证的路由
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(cfg))
	{
		// 用户相关路由
		protected.GET("/auth/profile", authHandler.GetProfile)
		protected.POST("/auth/logout", authHandler.Logout)

		// 管理员路由
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			// 这里可以添加管理员专用的路由
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Admin users endpoint"})
			})
		}

		// 连接管理路由
		connections := protected.Group("/connections")
		{
			connections.POST("", connectionHandler.CreateConnection)
			connections.GET("", connectionHandler.ListConnections)
			connections.GET("/:id", connectionHandler.GetConnection)
			connections.PUT("/:id", connectionHandler.UpdateConnection)
			connections.DELETE("/:id", connectionHandler.DeleteConnection)
			connections.POST("/:id/test", connectionHandler.TestConnection)

			// KV 管理路由
			connections.GET("/:id/kv", kvHandler.ListKeys)
			connections.GET("/:id/kv/*key", kvHandler.GetValue)
			connections.PUT("/:id/kv/*key", kvHandler.SetValue)
			connections.DELETE("/:id/kv/*key", kvHandler.DeleteKey)

			// 备份与导入路由
			connections.GET("/:id/backup/export", backupHandler.ExportBackup)
			connections.POST("/:id/backup/import", backupHandler.ImportBackup)
		}

		// KV 传输路由
		transferGroup := protected.Group("/transfer")
		{
			transferGroup.POST("", transferHandler.TransferKV)
			transferGroup.POST("/copy/:key", transferHandler.CopyKey)
		}
	}
}
