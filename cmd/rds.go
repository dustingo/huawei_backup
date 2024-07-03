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

// verifyRdsCmd 执行校验rds备份状态
var verifyRdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "verify rds backup status",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// 删除数据库备份
var deleteRdsBackupCmd = &cobra.Command{
	Use:   "rds",
	Short: "delete rds backup",
	Run: func(cmd *cobra.Command, args []string) {
		b := service.RdsDeleteBackup{Name: "rds delete"}
		b.Delete()
	},
}

func init() {
	backupCmd.AddCommand(rdsCmd)
	verifyCmd.AddCommand(verifyRdsCmd)
	deleteCmd.AddCommand(deleteRdsBackupCmd)
}
