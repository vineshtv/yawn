/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yawn",
	Short: "yet another wacky notes tool",
	Long:  `yawn - Yet Another Wacky Notes tool`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the yawn version",
	Long:  `Even yawn has versions. This command will reveal it.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.0")
	},
}

func init() {
	// Add version
	rootCmd.AddCommand(versionCmd)
}
