package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"etcd-admin-backend/internal/config"
	"etcd-admin-backend/internal/models"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.Config) error {
	var db *gorm.DB
	var err error
	if cfg.Database.Type == "sqlite" {
		db, err = gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Database,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// AutoMigrate 新模型
	if err := db.AutoMigrate(&models.User{}, &models.Connection{}, &models.KVItem{}); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("Database connected and migrated successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
