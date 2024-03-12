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

// syncStatsCmd represents the syncStats command
var syncStatsCmd = &cobra.Command{
	Use:   "syncStats",
	Short: "Fetch and save statistic based on 'games' table in db.",
	Long:  `Fetch and save statistic based on 'games' table in db.`,
	Run:   syncStats,
}

func syncStats(cmd *cobra.Command, args []string) {
	sp, err := controller.New()
	if err != nil {
		fmt.Println("controller error:", err)
		return
	}

	var queue = make(chan []balldontlie.Stats)
	process := func() {
		for qe := range queue {
			err := sp.Datastore.BulkSetStats(qe)
			if err != nil {
				fmt.Println("save error:", err)
				return
			}
		}
	}

	for i := 0; i < 10; i++ {
		go process()
	}

	ids, err := sp.Datastore.GameIds()
	if err != nil {
		fmt.Println("get all teams err:", err)
		return
	}
	fmt.Println("found", len(ids), "games")
	for _, i := range ids {
		_, err := sp.Api.PublishStatsByGameId(i, queue)
		if err != nil {
			fmt.Println("api err:", err)
			continue
		}
	}
	fmt.Println("completed sync")
}

func init() {
	rootCmd.AddCommand(syncStatsCmd)

}
