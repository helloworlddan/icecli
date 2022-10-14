package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/lensesio/tableprinter"
)

const (
	tripEndpoint = "/tripInfo/trip"
)

type Envelope struct {
	Trip Trip `json:"trip"`
}

type Trip struct {
	TrainType   string   `json:"trainType"`
	TrainNumber string   `json:"vzn"`
	StopInfo    StopInfo `json:"stopInfo"`
	Stops       []Stop   `json:"stops"`

	Train              string `header:"Train"`
	YourDestination    string `header:"Your Destination"`
	FinalDestination   string `header:"Final Destination"`
	NextStop           string `header:"Next Stop"`
	StopsToDestination int
	Progress           string `header:"Progress"`
}

type StopInfo struct {
	ScheduledNext     string `json:"scheduledNext"`
	ActualNext        string `json:"actualNext"`
	ActualLast        string `json:"actualLast"`
	ActualLastStarted string `json:"actualLastStarted"`
	FinalStationName  string `header:"FD" json:"finalStationName"`
}

type Stop struct {
	Station      Station       `json:"station"`
	TimeTable    TimeTable     `json:"timetable"`
	Track        Track         `json:"track"`
	Info         Info          `json:"info"`
	DelayReasons []DelayReason `json:"delayReasons"`
}

type Station struct {
	ID             string         `json:"evaNr"`
	Name           string         `json:"name"`
	GeoCoordinates GeoCoordinates `json:"geocoordinates"`
}

type GeoCoordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type TimeTable struct {
	ScheduledArrivalTime    uint64 `json:"scheduledArrivalTime"`
	ActualArrivalTime       uint64 `json:"actualArrivalTime"`
	ShowActualArrivalTime   bool   `json:"showActualArrivalTime"`
	ArrivalDelay            string `json:"arrivalDelay"`
	ScheduledDepartureTime  uint64 `json:"scheduledDepartureTime"`
	ActualDepartureTime     uint64 `json:"actualDepartureTime"`
	ShowActualDepartureTime bool   `json:"showActualDepartureTime"`
	DepartureDelay          string `json:"departureDelay"`
}

type Track struct {
	Scheduled string `json:"scheduled"`
	Actual    string `json:"actual"`
}

type Info struct {
	Status            int    `json:"status"`
	Passed            bool   `json:"passed"`
	PositionStatus    string `json:"positionStatus"`
	Distance          int    `json:"distance"`
	DistanceFromStart int    `json:"distanceFromStart"`
}

type DelayReason struct {
}

var (
	TripDestinationOverride string
	TripFilter              string
	TripOutput              string
)

var tripCmd = &cobra.Command{
	Use:   "trip",
	Short: "Print train trip info",
	Long:  `Print train trip info`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := refreshTrip()
		if err != nil {
			fail(err)
		}

		if StatusFilter != "" {
			switch StatusFilter {

			default:
				fail(errors.New("unknown filter field"))
			}

			fmt.Fprintf(os.Stderr, "\n")
			return
		}

		if StatusOutput == "table" {
			printer := tableprinter.New(os.Stdout)
			items := []Trip{t}
			printer.Print(items)
			return
		}
		if StatusOutput == "csv" {
			//fmt.Printf("%s,%s,%s,%f,%f,%f,%s\n", s)
			return
		}

		fail(errors.New("unrecognized output format"))
	},
}

func refreshTrip() (Trip, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", baseURL, tripEndpoint))
	if err != nil {
		return Trip{}, nil
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Trip{}, err
	}

	var envelope Envelope
	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return Trip{}, err
	}

	envelope.Trip.Train = fmt.Sprintf("%s%s", envelope.Trip.TrainType, envelope.Trip.TrainNumber)
	envelope.Trip.FinalDestination = envelope.Trip.StopInfo.FinalStationName
	envelope.Trip.YourDestination = func() string {
		if TripDestinationOverride == "" {
			return envelope.Trip.StopInfo.FinalStationName
		}
		for _, stop := range envelope.Trip.Stops {
			if stop.Station.Name == TripDestinationOverride {
				return stop.Station.Name
			}
		}
		fail(errors.New("overridden destination not found in schedule"))
		return ""
	}()
	envelope.Trip.StopsToDestination = func() int {
		numStops := 0
		for _, stop := range envelope.Trip.Stops {
			numStops++
			if stop.Station.Name == envelope.Trip.YourDestination {
				break
			}
		}
		return numStops
	}()
	envelope.Trip.NextStop = func() string {
		for _, stop := range envelope.Trip.Stops {
			if !stop.Info.Passed {
				return stop.Station.Name
			}
		}
		return "-"
	}()
	envelope.Trip.Progress = func() string {
		departedStops := 0
		for _, stop := range envelope.Trip.Stops {
			if stop.Info.Passed {
				departedStops++
			}
		}
		return fmt.Sprintf("%d/%d", departedStops, envelope.Trip.StopsToDestination)
	}()

	return envelope.Trip, nil
}

func init() {
	tripCmd.Flags().StringVarP(&TripDestinationOverride, "destination", "d", "", "Override for your destination")
	tripCmd.Flags().StringVarP(&TripOutput, "output", "o", "table", "Output format: table or csv")
	tripCmd.Flags().StringVarP(&TripFilter, "filter", "f", "", "Filter available fields")
	rootCmd.AddCommand(tripCmd)
}
