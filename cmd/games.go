package cmd

import (
	"fmt"
	"sync"

	"github.com/rajwalgautam/nba-client/pkg/balldontlie"
	"github.com/rajwalgautam/nba-client/pkg/controller"
	"github.com/spf13/cobra"
	"go.uber.org/ratelimit"
)

// gamesCmd represents the games command
var gamesCmd = &cobra.Command{
	Use:   "games",
	Short: "Fetch games with different parameters",
	Long:  `Fetch games with different parameters`,
	Run:   games,
}

func games(cmd *cobra.Command, args []string) {
	sp, err := controller.New()
	if err != nil {
		fmt.Println("controller error:", err)
		return
	}

	queue := make(chan balldontlie.Game, 0)
	var wg sync.WaitGroup
	rl := ratelimit.New(100)

	save := func(g balldontlie.Game) {
		defer wg.Done()
		defer rl.Take()
		err := sp.Datastore.SetGame(g)
		if err != nil {
			fmt.Println("save error:", err)
		}
	}

	process := func() {
		for qe := range queue {
			err := sp.Datastore.SetGame(qe)
			if err != nil {
				fmt.Println("save error:", err)
			}
		}
	}
	go process()

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

		for _, g := range games {
			queue <- g
		}
		fmt.Printf("saved %d games\n", len(games))
	}
}

func init() {
	rootCmd.AddCommand(gamesCmd)
}
