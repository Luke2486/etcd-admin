package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"etcd-admin-backend/internal/config"
)

// Migrate 执行数据库迁移
func Migrate(cfg *config.Config) error {
	// 初始化数据库连接用于迁移
	if err := InitDatabase(cfg); err != nil {
		return fmt.Errorf("failed to initialize database for migration: %w", err)
	}

	// SQLite使用GORM AutoMigrate就足够了
	if cfg.Database.Type == "sqlite" {
		log.Println("Database migration completed successfully (SQLite uses AutoMigrate)")
		return nil
	}

	// MySQL使用golang-migrate
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 创建迁移驱动
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// 创建迁移实例
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// 执行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// MigrateDown 回滚数据库迁移
func MigrateDown(cfg *config.Config) error {
	// SQLite不支持迁移回滚，直接删除数据库文件
	if cfg.Database.Type == "sqlite" {
		log.Println("SQLite migration rollback: please manually delete the database file if needed")
		return nil
	}

	// 初始化数据库连接用于迁移
	if err := InitDatabase(cfg); err != nil {
		return fmt.Errorf("failed to initialize database for migration: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 创建迁移驱动
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// 创建迁移实例
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// 回滚迁移
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Database migration rollback completed successfully")
	return nil
}
