package main

import (
	"os"
	"shuttle/databases"
	_ "github.com/lib/pq"
    "github.com/fatih/color"
	"github.com/pressly/goose/v3"
)

func main() {
	color.Yellow("Connecting to Database...")

	db, err := databases.PostgresConnection()
	if err != nil {
		color.Red("Failed to connect to PostgreSQL:", err)
	}

	sqlDB := db.DB

	err = os.Chdir("databases")
	if err != nil {
		color.Red("Failed to change directory to databases")
	}

	color.Yellow("Running migrations...")

	err = goose.Up(sqlDB, "./migrations")
	if err != nil {
		color.Red("Failed to run migration:", err)
	}

	color.Green("Migration successful")
}