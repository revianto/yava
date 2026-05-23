package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/revianto/yava/api/helpers"
)

const migrationsDir = "database/migrations"

func main() {
	_ = helpers.GetEnv("DB_HOST")

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [create|up|down|status|version|fix] [arguments]")
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "create" {
		if len(os.Args) < 3 {
			log.Fatal("Usage: create <migration_name>")
		}
		if err := goose.Create(nil, migrationsDir, os.Args[2], "sql"); err != nil {
			log.Fatal(err)
		}
		return
	}

	db := openDB()
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if _, err := goose.EnsureDBVersion(db); err != nil {
		log.Fatal("Could not initialize goose version table:", err)
	}

	var err error
	switch command {
	case "up":
		err = goose.Up(db, migrationsDir, goose.WithAllowMissing())
	case "down":
		err = goose.Down(db, migrationsDir)
	case "status":
		err = goose.Status(db, migrationsDir)
	case "version":
		err = goose.Version(db, migrationsDir)
	case "fix":
		err = goose.Fix(migrationsDir)
	case "reset":
		err = goose.Reset(db, migrationsDir)
	default:
		log.Fatal("Unknown command:", command)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func openDB() *sql.DB {
	sslmode := helpers.GetEnv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		helpers.GetEnv("DB_HOST"),
		helpers.GetEnv("DB_PORT"),
		helpers.GetEnv("DB_USERNAME"),
		helpers.GetEnv("DB_PASSWORD"),
		helpers.GetEnv("DB_DATABASE"),
		sslmode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Could not ping database:", err)
	}
	return db
}
