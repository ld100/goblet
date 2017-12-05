package main

import (
	"os"

	"github.com/ld100/goblet/persistence"
	"github.com/ld100/goblet/persistence/migrate"
	"github.com/ld100/goblet/server/rest"
)

func main() {
	// Prepare initial data: create db, run migrations and seeds
	prepareData()

	// Launch CHI-based RESTful HTTP server
	rest.Serve()
}

// TODO: Move whole chain of commands to persistence package
func prepareData() {
	// Create database if not exist
	ds := &persistence.DataSource{}
	// Fetch database credentials from ENVIRONMENT
	ds.FetchENV()
	ds.CreateDB(os.Getenv("DB_NAME"))

	// Initiate global ORM var
	persistence.InitGormDB(ds)

	// Run migrations
	migrate.Migrate()

	// Run db seed
	migrate.Seed()
}
