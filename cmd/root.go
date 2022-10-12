package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	baseURL = "https://iceportal.de/api1/rs"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "icecli",
	Short: "icecli is a tool to query on-train ICE API endpoints for real-time information about train and journey",
	Long:  `icecli is a tool to query on-train ICE API endpoints for real-time information about train and journey`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
