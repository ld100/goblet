package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goblet",
	Short: "Goblet Server",
	Long:  `Goblet is an authentication server`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("Print: " + strings.Join(args, " "))
		fmt.Println("No commands given. Run 'goblet --help' for usage help.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
