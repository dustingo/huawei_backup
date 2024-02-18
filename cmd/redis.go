/*
Copyright © 2024 zack
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"huawei_backup/service"
	"time"
)

// redisCmd represents the redis command
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "do redis backup",
	Run: func(cmd *cobra.Command, args []string) {
		var db service.Database
		db = service.RedisBackup{
			Name: fmt.Sprintf("redis-manualbackup-%v", time.Now().Format("20060102150405")),
		}
		db.Backup()
	},
}

func init() {
	backupCmd.AddCommand(redisCmd)
}
