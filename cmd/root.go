package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	baseURL = "https://iceportal.de/api1/rs"
)

var (
	Output string
	Filter string
)

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

func init() {
	tripCmd.PersistentFlags().StringVarP(&Output, "output", "o", "table", "Output format: table or csv")
	tripCmd.PersistentFlags().StringVarP(&Filter, "filter", "f", "", "Filter available fields")
}

func unixMillisToTime(millis uint64) time.Time {
	seconds := int64(millis / 1000)
	return time.Unix(seconds, 0)
}

func formatTimeDelta(d time.Duration) string {
	return strings.Split(d.String(), ".")[0]
}
