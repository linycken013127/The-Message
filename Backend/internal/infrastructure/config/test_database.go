package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// InitTestDB 初始化測試資料庫
func InitTestDB() *gorm.DB {
	_, err := LoadEnvPath()
	if err != nil {
		return nil
	}

	db := NewDatabase()
	return db
}

// LoadEnvPath 載入環境變數路徑
func LoadEnvPath() (string, error) {
	envPath, err := getEnvPath()

	err = godotenv.Load(envPath)
	if err != nil {
		return "", err
	}

	return envPath, err
}

// getEnvPath 取得環境變數路徑
func getEnvPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	envPath := cwd + "/../../.env"
	return envPath, err
}

// BaseTestDSN 取得測試用 DSN
func BaseTestDSN() string {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	return dsn
}
