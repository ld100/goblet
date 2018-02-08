package cmd

import (
	"fmt"
	"strings"

	"github.com/ld100/goblet/pkg/persistence/setup"
	"github.com/ld100/goblet/pkg/util/config"
	"github.com/ld100/goblet/pkg/util/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbCreateCmd)
}

var dbCreateCmd = &cobra.Command{
	Use:   "dbcreate",
	Short: "Create initial PostgreSQL database for Goblet",
	Long:  `Create initial PostgreSQL database for Goblet`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating database with arguments: " + strings.Join(args, " "))
		dbName := ""
		if len(args) > 0 {
			dbName = args[0]
		}

		cfg := &config.Config{}
		log := logger.New(cfg)
		err := setup.CreateDB(cfg, dbName)
		if err != nil {
			log.Fatal("cannot create database database ", err)
		} else {
			fmt.Println("database created succesfully")
		}
	},
}
