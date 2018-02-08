package cmd

import (
	"fmt"

	"github.com/ld100/goblet/pkg/persistence/seed"
	"github.com/ld100/goblet/pkg/server/env"
	"github.com/ld100/goblet/pkg/util/config"
	"github.com/ld100/goblet/pkg/util/logger"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "dbseed",
	Short: "Run database seeds",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running database seed...")

		// Build initial env object, which would be passed to HTTP server
		env := &env.Env{}

		// Create global configuration object and attach it to environment
		cfg := &config.Config{}
		env.Config = cfg

		// Create global logger object and pass it to environment
		log := logger.New(cfg)
		env.Logger = log

		err := seed.SeedDB(env)
		if err != nil {
			log.Fatal("cannot seed database", err)
		} else {
			fmt.Println("database seeds ran succesfully")
		}
	},
}
