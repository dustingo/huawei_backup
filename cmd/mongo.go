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

// mongoCmd represents the mongo command
var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "do mongodb bakcup",

	Run: func(cmd *cobra.Command, args []string) {
		var db service.Database
		db = service.MongoBackup{
			Name: fmt.Sprintf("mongodb-manualbackup-%v", time.Now().Format("20060102150405")),
		}
		db.Backup()
	},
}

func init() {
	backupCmd.AddCommand(mongoCmd)
}
