/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mmitten/arbitage/theoddsinterface"
	"github.com/spf13/cobra"
)

type OddsFlags struct {
	Market string
}

var odds = OddsFlags{}

// oddsCmd represents the odds command
var oddsCmd = &cobra.Command{
	Use:   "odds",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sports, errSports := theoddsinterface.GetSports()
		if errSports != nil {
			cmd.PrintErrln(errSports)
		}
		results, errResults := sports.CheckMarket(odds.Market)
		if errResults != nil {
			cmd.PrintErrln(errResults)
		}
		for _, result := range results {
			cmd.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(oddsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// oddsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// oddsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	oddsCmd.Flags().StringVar(&odds.Market, "market", "", "Set the market to check.")
}
