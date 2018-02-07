package main

import (
	"github.com/ld100/goblet/pkg/persistence/setup"
	"github.com/ld100/goblet/pkg/server/rest"
)

func main() {
	// Prepare initial data: create db, run migrations and seeds
	// TODO: Change _ with dba and pass db to rest.Serve()
	_, err := setup.SetupDatabases()
	if err != nil {
	//	Put error handling here
	}

	// Launch CHI-based RESTful HTTP server
	rest.Serve()
}
