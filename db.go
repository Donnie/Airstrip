package main

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initSQLDB() *gorm.DB {
	// init DB
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE")), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database")
		os.Exit(1)
	}
	fmt.Println("Successfully connected to database")

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(10)
	migrateSQLUp()

	return db
}
