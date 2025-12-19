package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDataBase() *gorm.DB {
	host := os.Getenv("SQL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable&TrustServerCertificate=true",
		os.Getenv("SQL_USER"),
		os.Getenv("SQL_PASSWORD"),
		host,
		os.Getenv("SQL_PORT"),
		os.Getenv("SQL_DB"),
	)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		DryRun:                 false,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get generic database object: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("Database connection successful!")
	return db
}
