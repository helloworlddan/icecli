package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/lensesio/tableprinter"
)

type StopView struct {
	Station           string `header:"Station"`
	Track             string `header:"Track"`
	TimeToArrival     string `header:"Arriving"`
	RemainingDistance string `header:"Remaining Distance"`
	DelayReasons      string `header:"Reasons for delay"`
}

var stopsCmd = &cobra.Command{
	Use:   "stops",
	Short: "Print trip stop info",
	Long:  `Print trip stop info`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := refreshTrip()
		if err != nil {
			fail(err)
		}

		stopViews := []StopView{}
		for idx, stop := range t.Stops {
			stopView := StopView{}
			stopView.Station = stop.Station.Name
			stopView.Track = stop.Track.Actual
			stopView.TimeToArrival = func() string {
				arrivalTime := unixMillisToTime(stop.TimeTable.ActualArrivalTime)
				remainingDuration := time.Until(arrivalTime)
				if idx == 0 {
					return ""
				}
				if remainingDuration < 0 {
					return "-"
				}
				if stop.Station.Name == t.YourDestination && remainingDuration < time.Duration(time.Minute*3) {
					return "GET OUT NOW"
				}
				return formatTimeDelta(remainingDuration)
			}()
			stopView.RemainingDistance = func() string {
				if idx == 0 {
					return ""
				}
				// Is this actually correct or should ActualPosition be replaced with info from the last stop?
				traveledDistance := t.ActualPosition + t.DistanceFromLastStop
				remainingMeters := stop.Info.DistanceFromStart - traveledDistance
				roundedKilometers := remainingMeters / 1000
				if roundedKilometers <= 0 {
					return "-"
				}

				return fmt.Sprintf("%d km", roundedKilometers)
			}()
			stopView.DelayReasons = func() string {
				if idx == 0 {
					return ""
				}
				reasons := []string{}
				for _, r := range stop.DelayReasons {
					reasons = append(reasons, r.Message)
				}

				concatReasons := strings.Join(reasons, "; ")
				if concatReasons == "" {
					return "-"
				}
				return concatReasons
			}()
			stopViews = append(stopViews, stopView)
			if stop.Station.Name == t.YourDestination {
				break
			}
		}

		if Filter != "" {
			switch Filter {
			case "STATION":
				stations := []string{}
				for _, stopView := range stopViews {
					stations = append(stations, stopView.Station)
				}
				fmt.Println(strings.Join(stations, ","))
			case "TRACK":
				tracks := []string{}
				for _, stopView := range stopViews {
					tracks = append(tracks, stopView.Track)
				}
				fmt.Println(strings.Join(tracks, ","))
			case "ARRIVING":
				arrivalDurations := []string{}
				for _, stopView := range stopViews {
					arrivalDurations = append(arrivalDurations, stopView.Track)
				}
				fmt.Println(strings.Join(arrivalDurations, ","))
			case "REMAINING DISTANCE":
				distances := []string{}
				for _, stopView := range stopViews {
					distances = append(distances, stopView.RemainingDistance)
				}
				fmt.Println(strings.Join(distances, ","))
			case "REASONS FOR DELAY":
				delayReasons := []string{}
				for _, stopView := range stopViews {
					delayReasons = append(delayReasons, stopView.Track)
				}
				fmt.Println(strings.Join(delayReasons, ","))
			default:
				fail(errors.New("unknown filter field"))
			}

			fmt.Fprintf(os.Stderr, "\n")
			return
		}

		if Output == "table" {
			printer := tableprinter.New(os.Stdout)
			printer.Print(stopViews)
			return
		}
		if Output == "csv" {
			for _, stopView := range stopViews {
				fmt.Printf("%s,%s,%s,%s,%s\n", stopView.Station, stopView.Track, stopView.TimeToArrival, stopView.RemainingDistance, stopView.DelayReasons)
			}
			return
		}

		fail(errors.New("unrecognized output format"))
	},
}

func init() {
	tripCmd.AddCommand(stopsCmd)
}
