package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "0.0.4"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information.",
	Long:  `Print version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("icecli v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
