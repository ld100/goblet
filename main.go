package main

import (
	"fmt"
	"os"
	"github.com/ld100/goblet/util/database"
	//"github.com/ld100/goblet/users"
	"github.com/ld100/goblet/util/environment"
	"github.com/ld100/goblet/util/migrate"
	"github.com/ld100/goblet/server/rest"
)

func main() {
	// Prepare initial data: create db, run migrations and seeds
	prepareData()

	// Launch CHI-based RESTful HTTP server
	rest.Serve()
}

func prepareData() {
	// Create database if not exist
	database.CreateDB(os.Getenv("DB_NAME"))

	// Initiate global ORM var
	connString := fmt.Sprintf(
		"host=%v user=%v dbname=%v sslmode=disable password=%v",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	environment.InitGDB(connString)

	// Run migrations
	migrate.Migrate()

	// Run db seed
	migrate.Seed()
}
