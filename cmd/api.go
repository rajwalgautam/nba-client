package cmd

import (
	"fmt"

	"github.com/rajwalgautam/nba-client/pkg/controller"
	"github.com/rajwalgautam/nba-client/pkg/db"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "validate api key",
	Long:  `Validate api key for balldontlie api`,
	Run:   api,
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.PersistentFlags().String("key", "", "set api key")
}

func api(cmd *cobra.Command, args []string) {
	key, set := getFlag(cmd, "key")
	if set {
		c, err := db.New()
		if err != nil {
			fmt.Println("db error", err)
			return
		}
		err = c.SetApiKey(key)
		if err != nil {
			fmt.Println("set key error:", err)
			return
		}
		fmt.Println("key set")
		return
	}

	sp, err := controller.New()
	if err != nil {
		fmt.Println("stats processor error:", err)
	}
	err = sp.Api.Ping()
	if err != nil {
		fmt.Println("check balldontlie error:", err)
		return
	}
	fmt.Println("balldontlie api is up")
}
