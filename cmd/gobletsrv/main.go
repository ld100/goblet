package main

import (
	"github.com/ld100/goblet/pkg/persistence/setup"
	"github.com/ld100/goblet/pkg/server/env"
	"github.com/ld100/goblet/pkg/server/rest"
	"github.com/ld100/goblet/pkg/util/config"
)

func main() {
	// Build initial env object, which would be passed to HTTP server
	env := &env.Env{}

	// Create global configuration object and attach it to environment
	cfg := &config.Config{}
	env.Config = cfg

	// TODO: Put global logger to Env

	// Prepare initial data: create db, run migrations and seeds
	// Take appropriate configuration data from cfg object
	db, err := setup.SetupDatabases(cfg)
	if err != nil {
		//	Put error handling here
	}
	env.DB = db

	// Launch CHI-based RESTful HTTP server
	rest.Serve(env)
}
