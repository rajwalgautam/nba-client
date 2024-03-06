package cmd

import "github.com/spf13/cobra"

func getFlag(cmd *cobra.Command, fl string) (string, bool) {
	isSet := cmd.Flags().Lookup(fl).Changed
	s, err := cmd.Flags().GetString(fl)
	if err != nil {
		return "", false
	}
	return s, isSet
}
