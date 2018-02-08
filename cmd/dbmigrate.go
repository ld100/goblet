package cmd

import (
	"fmt"

	"github.com/ld100/goblet/pkg/persistence/migrate"
	"github.com/ld100/goblet/pkg/server/env"
	"github.com/ld100/goblet/pkg/util/config"
	"github.com/ld100/goblet/pkg/util/logger"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "dbmigrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running database migrations...")

		// Build initial env object, which would be passed to HTTP server
		env := &env.Env{}

		// Create global configuration object and attach it to environment
		cfg := &config.Config{}
		env.Config = cfg

		// Create global logger object and pass it to environment
		log := logger.New(cfg)
		env.Logger = log

		err := migrate.MigrateDB(env)
		if err != nil {
			log.Fatal("cannot migrate database", err)
		} else {
			fmt.Println("database migrated succesfully")
		}
	},
}
