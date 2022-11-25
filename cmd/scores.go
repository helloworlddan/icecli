package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/lensesio/tableprinter"
	"github.com/spf13/cobra"
)

const (
	scoresEndpoint = "/ticker/live-or-past/wm2022"
)

type Score struct {
	Competition         string  `header:"Competition" json:"competition"`
	Venue               string  `header:"Venue" json:"venue"`
	Team1               Team    `json:"team1"`
	Goals1              int     `json:"goals1"`
	Team2               Team    `json:"team2"`
	Goals2              int     `json:"goals2"`
	Match               string  `header:"Match"`
	Goals               string  `header:"Goals"`
	Events              []Event `json:"events"`
	MostRecentEventText string  `header:"Last Event"`
}

type Team struct {
	Name string `header:"Name" json:"name"`
}

type Event struct {
	Text string `header:"Event" json:"eventText"`
}

var scoresCmd = &cobra.Command{
	Use:   "scores",
	Short: "Print sports tournament scores info",
	Long:  `Print sports tournament scores info`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := refreshScores()
		if err != nil {
			fail(err)
		}

		if Output == "table" {
			printer := tableprinter.New(os.Stdout)
			items := s
			printer.Print(items)
			return
		}
		if Output == "csv" {
			for _, score := range s {
				fmt.Printf("%s %s %s %s %s\n", score.Competition, score.Venue, score.Match, score.Goals, score.Events[0].Text)
			}
			return
		}

		fail(errors.New("unrecognized output format"))
	},
}

func refreshScores() ([]*Score, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", baseURL, scoresEndpoint))
	if err != nil {
		return []*Score{}, nil
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []*Score{}, err
	}

	var scores []*Score
	err = json.Unmarshal(data, &scores)
	if err != nil {
		return []*Score{}, err
	}

	for _, score := range scores {
		score.Match = fmt.Sprintf("%s vs. %s", score.Team1.Name, score.Team2.Name)
		score.Goals = fmt.Sprintf("%d:%d", score.Goals1, score.Goals2)
		score.MostRecentEventText = score.Events[0].Text
	}

	return scores, nil
}

func init() {
	rootCmd.AddCommand(scoresCmd)
}
