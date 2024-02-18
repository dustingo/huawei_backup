/*
Copyright © 2024 zack
*/
package cmd

import (
	"fmt"
	"huawei_backup/service"
	"time"

	"github.com/spf13/cobra"
)

// rdsCmd represents the rds command
var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "do rds backup",
	Run: func(cmd *cobra.Command, args []string) {
		var db service.Database
		db = service.RdsBackup{
			Name: fmt.Sprintf("rds-manualbackup-%v", time.Now().Format("20060102150405")),
		}
		db.Backup()
	},
}

func init() {
	backupCmd.AddCommand(rdsCmd)
}
