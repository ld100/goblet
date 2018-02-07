package main

import (
	"github.com/ld100/goblet/pkg/persistence/setup"
	"github.com/ld100/goblet/pkg/server/rest"
	"github.com/ld100/goblet/pkg/server/env"
)

func main() {
	// Build initial env object, which would be passed to HTTP server
	env := &env.Env{}
	// TODO: Put global configuration to Env
	// TODO: Put global logger to Env

	// Prepare initial data: create db, run migrations and seeds
	db, err := setup.SetupDatabases()
	if err != nil {
		//	Put error handling here
	}
	env.DB = db

	// Launch CHI-based RESTful HTTP server
	rest.Serve(env)
}
