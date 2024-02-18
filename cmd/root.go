/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "huawei_backup",
	Short:         "huawei cloud database backup",
	SilenceUsage:  false,
	SilenceErrors: false,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config := viper.GetString("config")
		if config == "" {
			fmt.Println("config file is empty")
			os.Exit(-1)
		}
		viper.SetConfigFile(config)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		viper.Set("push", viper.GetBool("push"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	rootCmd.PersistentFlags().BoolP("push", "p", false, "push to prometheus pushgateway")
	rootCmd.PersistentFlags().String("config", "", "config file")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := viper.BindPFlag("push", rootCmd.PersistentFlags().Lookup("push")); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
