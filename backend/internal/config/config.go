package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Type     string
	Path     string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ServerConfig struct {
	Port    string
	JWTKey  string
	GinMode string
}

func LoadConfig() *Config {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables and defaults")
	}

	return &Config{
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "mysql"),
			Path:     getEnv("DB_PATH", "etcd-admin.db"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_DATABASE", "etcd_admin"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			JWTKey:  getEnv("JWT_SECRET", "your-secret-key"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
