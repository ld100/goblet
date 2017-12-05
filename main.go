package main

import (
	"github.com/ld100/goblet/persistence/setup"
	"github.com/ld100/goblet/server/rest"
)

func main() {
	// Prepare initial data: create db, run migrations and seeds
	setup.SetupDatabases()

	// Launch CHI-based RESTful HTTP server
	rest.Serve()
}
