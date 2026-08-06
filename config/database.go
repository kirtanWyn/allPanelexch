package config

import (
	"database/sql"
	"fmt"
	// "time"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase() error {

// fmt.Println("DB_HOST :", os.Getenv("DB_HOST"))
// fmt.Println("DB_PORT :", os.Getenv("DB_PORT"))
// fmt.Println("DB_USER :", os.Getenv("DB_USER"))
// fmt.Println("DB_NAME :", os.Getenv("DB_NAME"))
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	// Connection Pool Settings
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(30 * time.Minute)
	// db.SetConnMaxIdleTime(10 * time.Minute)

	// Verify Connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	DB = db

	fmt.Println("✅ MySQL Connected Successfully")

	return nil
}

func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
