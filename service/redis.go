package service

import (
	dcs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dcs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dcs/v2/model"
	"github.com/spf13/viper"
	"huawei_backup/logger"
	"huawei_backup/pkg/client"
	"huawei_backup/pkg/push"
	"strings"
	"sync"
)

type RedisBackup struct {
	Name string
}

func (r RedisBackup) Backup() {
	var wg sync.WaitGroup
	var backupFormatBackupInstanceBody model.BackupInstanceBodyBackupFormat
	dcsClient := client.NewRedisClient()
	if strings.ToLower(viper.GetString("redis.format")) == "aof" {
		backupFormatBackupInstanceBody = model.GetBackupInstanceBodyBackupFormatEnum().AOF
	} else {
		backupFormatBackupInstanceBody = model.GetBackupInstanceBodyBackupFormatEnum().RDB
	}
	requests := make([]*model.CopyInstanceRequest, 0)
	for _, instance := range viper.GetStringSlice("redis.instanceId") {
		req := &model.CopyInstanceRequest{}
		body := &model.BackupInstanceBody{
			BackupFormat: &backupFormatBackupInstanceBody,
			Remark:       &r.Name,
		}
		req.InstanceId = instance
		req.Body = body
		requests = append(requests, req)
	}
	for _, request := range requests {
		wg.Add(1)
		go func(client *dcs.DcsClient, req *model.CopyInstanceRequest) {
			defer wg.Done()
			response, err := client.CopyInstance(req)
			if err != nil {
				logger.Logger.Errorw("create redis manual backup error", "error", err.Error())
				return
			}
			push.Push("redis", req.InstanceId)
			logger.Logger.Infow("create redis manual backup success", "InstanceId", req.InstanceId, "BackupId", response.BackupId)
		}(dcsClient, request)
	}
	wg.Wait()
}
