/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/rajwalgautam/nba-client/pkg/db"
	"github.com/spf13/cobra"
)

// initdbCmd represents the initdb command
var initdbCmd = &cobra.Command{
	Use:   "initdb",
	Short: "Init DB",
	Long:  `Migrate tables`,
	Run:   initdb,
}

func init() {
	rootCmd.AddCommand(initdbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initdbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initdbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initdb(cmd *cobra.Command, args []string) {
	_, err := db.New()
	if err != nil {
		fmt.Println("db error:", err)
		return
	}
	fmt.Println("db created")
}
