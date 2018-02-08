package cmd

import (
	"fmt"

	"github.com/ld100/goblet/pkg/persistence/setup"
	"github.com/ld100/goblet/pkg/server/env"
	"github.com/ld100/goblet/pkg/server/rest"
	"github.com/ld100/goblet/pkg/util/config"
	"github.com/ld100/goblet/pkg/util/logger"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve HTTP requests",
	Long:  `Launches internal REST/gRPC servers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting HTTP server...")

		// Build initial env object, which would be passed to HTTP server
		env := &env.Env{}

		// Create global configuration object and attach it to environment
		cfg := &config.Config{}
		env.Config = cfg

		// Create global logger object and pass it to environment
		log := logger.New(cfg)
		env.Logger = log

		// Prepare initial data: create db, run migrations and seeds
		// Take appropriate configuration data from cfg object
		db, err := setup.SetupDatabases(cfg)
		if err != nil {
			log.Fatal("cannot set up the database ", err)
		} else {
			env.DB = db

			// Launch CHI-based RESTful HTTP server
			rest.Serve(env)
		}
	},
}
