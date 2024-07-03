// Package cmd
// Created by zack
// 校验数据最后一个备份的状态
package cmd

import "github.com/spf13/cobra"

var verifyCmd = &cobra.Command{
	Use:           "verify",
	Short:         "verify the last database backup status",
	SilenceErrors: false,
	SilenceUsage:  false,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
