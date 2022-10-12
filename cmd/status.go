package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/lensesio/tableprinter"
)

const (
	statusEndpoint = "/status"
)

type Status struct {
	TrainType  string  `header:"Train Type",json:"trainType"`
	WagonClass string  `header:"Wagon Class",json:"wagonClass"`
	Internet   string  `header:"Internet",json:"internet"`
	Speed      float32 `header:"Speed",json:"speed"`
	Latitude   float32 `header:"Latitude",json:"latitude"`
	Longitude  float32 `header:"Longitude",json:"longitude"`
	GPSStatus  string  `header:"GPS",json:"gpsStatus"`

	TimeMillis uint64 `json:"serverTime"`
}

var (
	Filter string
	Output string
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print train status info",
	Long:  `Print train status info`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := refreshStatus()
		if err != nil {
			fail(err)
		}

		if Filter != "" {
			switch Filter {
			case "TRAIN TYPE":
				fmt.Printf("%s", s.TrainType)
			case "WAGON CLASS":
				fmt.Printf("%s", s.WagonClass)
			case "INTERNET":
				fmt.Printf("%s", s.Internet)
			case "SPEED":
				fmt.Printf("%f", s.Speed)
			case "LATITUDE":
				fmt.Printf("%f", s.Latitude)
			case "Longitude":
				fmt.Printf("%f", s.Longitude)
			case "GPS":
				fmt.Printf("%s", s.GPSStatus)
			}

			fmt.Fprintf(os.Stderr, "\n")
			return
		}

		if Output == "table" {
			printer := tableprinter.New(os.Stdout)
			items := []Status{s}
			printer.Print(items)
			return
		}
		if Output == "csv" {
			fmt.Printf("%s,%s,%s,%f,%f,%f,%s\n", s.TrainType, s.WagonClass, s.Internet, s.Speed, s.Latitude, s.Longitude, s.GPSStatus)
			return
		}

		fail(errors.New("unrecognized output format"))
	},
}

func refreshStatus() (Status, error) {
	file, err := os.Open("./samples/status.json")
	if err != nil {
		return Status{}, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return Status{}, err
	}

	var status Status
	err = json.Unmarshal(data, &status)
	if err != nil {
		return Status{}, err
	}

	return status, nil
}

func init() {
	statusCmd.Flags().StringVarP(&Output, "output", "o", "table", "Output format: table or csv")
	statusCmd.Flags().StringVarP(&Filter, "filter", "f", "", "Filter available fields")
	rootCmd.AddCommand(statusCmd)
}
