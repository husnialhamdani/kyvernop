/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Run loads of resources creation",
	Long:  `Run loadsd of resources creation`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("execute called")
		scales, _ := cmd.Flags().GetString("scales")
		fmt.Println(scales)
	},
}

func init() {
	executeCmd.Flags().StringP("scales", "s", "xs", "choose the scale size (small/medium/large/xl) default: xs")
	rootCmd.AddCommand(executeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
