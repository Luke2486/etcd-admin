package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"etcd-admin-backend/internal/config"
	"etcd-admin-backend/internal/handlers"
	"etcd-admin-backend/internal/middleware"
	"etcd-admin-backend/internal/services"
	"etcd-admin-backend/pkg/database"
)

func main() {
	// 命令行参数
	var migrate = flag.Bool("migrate", false, "Run database migration")
	var rollback = flag.Bool("rollback", false, "Rollback database migration")
	flag.Parse()

	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.GinMode)

	// 处理迁移命令
	if *migrate {
		log.Println("Running database migration...")
		if err := database.Migrate(cfg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed successfully")
		return
	}

	if *rollback {
		log.Println("Rolling back database migration...")
		if err := database.MigrateDown(cfg); err != nil {
			log.Fatalf("Migration rollback failed: %v", err)
		}
		log.Println("Migration rollback completed successfully")
		return
	}

	// 初始化数据库连接
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化etcd服务
	etcdService := services.NewEtcdService()

	// 创建 Gin 路由
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())

	// 基本路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "etcd-admin backend server running",
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// 设置API路由
	handlers.SetupRoutes(r, cfg, etcdService)

	// 启动服务器
	port := ":" + cfg.Server.Port
	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s", cfg.Server.GinMode)
	log.Fatal(r.Run(port))
}
