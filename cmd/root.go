/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vineshtv/yawn/config"
)

var (
	// config file.
	cfgFile string
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "yawn",
		Short: "yet another wacky notes tool",
		Long:  `yawn - Yet Another Wacky Notes tool`,
	}
)

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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.yawn.toml)")
	// Add version
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use the config file passed in the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Read config file from default location

		// Find the home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Set config file in home directory
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".yawn")
	}

	// viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		err := viper.Unmarshal(&config.Config)
		if err != nil {
			fmt.Println("Error processing config file:", viper.ConfigFileUsed())
			os.Exit(1)
		}
	} else {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
}
