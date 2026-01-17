package config

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDatabase 建立資料庫連線
func NewDatabase() *gorm.DB {
	dsn := DefaultDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	return db
}

// DefaultDSN 取得預設 DSN
func DefaultDSN() string {
	dsn := BaseDSN()

	val := url.Values{}
	val.Add("parseTime", "true")
	val.Add("loc", "Local")

	dsn = fmt.Sprintf("%s?%s", dsn, val.Encode())
	return dsn
}

// BaseDSN 取得基礎 DSN
func BaseDSN() string {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	return dsn
}
