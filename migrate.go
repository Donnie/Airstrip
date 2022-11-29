package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func migrateSQLUp() {
	db, err := sql.Open("sqlite", os.Getenv("DB_FILE"))
	if err != nil {
		fmt.Println("Failed to open SQL file")
		os.Exit(1)
	}
	defer db.Close()

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		fmt.Println("Failed to open SQL file")
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "ql", driver)
	if err != nil {
		fmt.Println("Failed to migrate")
		os.Exit(1)
	}

	// migration upto 10 steps
	m.Steps(10)
}
