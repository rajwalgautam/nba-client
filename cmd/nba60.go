/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/rajwalgautam/nba-client/pkg/balldontlie"
	"github.com/rajwalgautam/nba-client/pkg/controller"
	"github.com/spf13/cobra"
)

// nba60Cmd represents the nba60 command
var nba60Cmd = &cobra.Command{
	Use:   "nba60",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: nba60,
}

func nba60(cmd *cobra.Command, args []string) {
	c, err := controller.New()
	if err != nil {
		fmt.Println("controller err:", err)
		return
	}

	ids, err := c.Datastore.PlayerIds()
	if err != nil {
		fmt.Println("player ids err:", err)
		return
	}
	if len(ids) == 0 {
		fmt.Println("no player ids")
		return
	}
	players := make([]balldontlie.Player, len(ids))
	for i, pid := range ids {
		pl, err := c.Datastore.GetPlayer(pid)
		if err != nil {
			fmt.Printf("db error for %d: %s\n", i, err)
			continue
		}
		players[i] = pl
		count := 0
		for _, g := range pl.Games {
			if sum, s := g.Sixty(); s {
				count++
				fmt.Printf("%s %s (%s) TOTAL %d - %d pts | %d reb | %d ast | %d blk | %d stl\n", pl.First, pl.Last, g.Date, sum, g.Points, g.Rebounds, g.Assists, g.Blocks, g.Steals)
			}
		}
		if count > 0 {
			fmt.Printf("%s %s nba60 total: %d\n\n", pl.First, pl.Last, count)
		}
	}
	// for _, p := range players {
	// 	count := 0
	// 	for _, g := range p.Games {
	// 		if sum, s := g.Sixty(); s {
	// 			count++
	// 			fmt.Printf("%s %s (%s) TOTAL %d - %d pts | %d reb | %d ast | %d blk | %d stl\n", p.First, p.Last, g.Date, sum, g.Points, g.Rebounds, g.Assists, g.Blocks, g.Steals)
	// 		}
	// 	}
	// 	if count > 0 {
	// 		fmt.Printf("%s %s nba60 total: %d\n\n", p.First, p.Last, count)
	// 	}
	// }
}

func init() {
	statsCmd.AddCommand(nba60Cmd)
}
