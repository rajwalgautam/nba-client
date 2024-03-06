/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sync"

	"github.com/rajwalgautam/nba-client/pkg/balldontlie"
	"github.com/rajwalgautam/nba-client/pkg/statsprocessor"
	"github.com/spf13/cobra"
)

// gamesCmd represents the games command
var gamesCmd = &cobra.Command{
	Use:   "games",
	Short: "Fetch games with different parameters",
	Long:  `Fetch games with different parameters`,
	Run:   games,
}

func games(cmd *cobra.Command, args []string) {
	sp, err := statsprocessor.New()
	if err != nil {
		fmt.Println("statsprocessor error:", err)
		return
	}

	var wg sync.WaitGroup

	save := func(g balldontlie.Game) {
		err := sp.Datastore.SetGame(g)
		if err != nil {
			fmt.Println("save error:", err)
		}
		gameStats, err := sp.Api.StatsByGameId(g.Id)
		if err != nil {
			fmt.Printf("get game stats error (%svs%s %s): %s\n", g.AwayTeam.Abbr, g.HomeTeam.Abbr, g.Date, err)
			wg.Done()
			return
		}
		err = sp.Datastore.BulkSetStats(gameStats)
		if err != nil {
			fmt.Println(err)
			return
		}

		wg.Done()
	}

	date, set := getFlag(cmd, "date")
	if set {
		games, err := sp.Api.GamesOnDate(date)
		if err != nil {
			fmt.Println("games on date error:", err)
			return
		}
		if len(games) == 0 {
			fmt.Println("no games")
			return
		}
		wg.Add(len(games))
		for _, g := range games {
			go save(g)
		}
		wg.Wait()
		fmt.Printf("saved %d games\n", len(games))
	}

	start, startSet := getFlag(cmd, "start")
	end, endSet := getFlag(cmd, "end")
	if startSet && endSet {
		games, err := sp.Api.GamesDateRange(start, end)
		if err != nil {
			fmt.Println("games date range error:", err)
			return
		}
		if len(games) == 0 {
			fmt.Println("no games")
			return
		}
		wg.Add(len(games))
		for _, g := range games {
			go save(g)
		}
		wg.Wait()
		fmt.Printf("saved %d games\n", len(games))
	}
}

func init() {
	rootCmd.AddCommand(gamesCmd)

	gamesCmd.PersistentFlags().StringP("date", "d", "", "get nba games at date")
	gamesCmd.PersistentFlags().StringP("start", "s", "", "start date for nba games, 'end' must also be set")
	gamesCmd.PersistentFlags().StringP("end", "e", "", "end date for nba games, 'start' must also be set")

}
