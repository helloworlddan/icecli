package cmd

import (
	"github.com/spf13/cobra"
)

const (
	statusEndpoint = "/status"
)

type status struct {
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print train status info",
	Long:  `Print train status info`,
	Run: func(cmd *cobra.Command, args []string) {

		// TODO do something

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
