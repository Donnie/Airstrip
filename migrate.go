package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func pgConnectionString() string {
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASS")
	dbas := os.Getenv("PG_DBAS")
	port := os.Getenv("PG_PORT")
	host := os.Getenv("PG_HOST")
	ssl := os.Getenv("PG_SSL")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbas, ssl)
}

func migrateUp() {
	m, err := migrate.New("file://migrations", pgConnectionString())
	if err != nil {
		panic(err)
	}
	// migration upto 10 steps
	m.Steps(10)
}

func migrateSqlUp() {
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
