package cmd

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:           "delete",
	Short:         "delete database backup",
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
	rootCmd.AddCommand(deleteCmd)
}
