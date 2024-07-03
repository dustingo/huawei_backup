package service

import (
	"fmt"
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
type RedisDeleteBackup struct {
	Name string
}

func (r RedisDeleteBackup) Delete() {
	var backupsList map[string][]model.BackupRecordResponse = make(map[string][]model.BackupRecordResponse)
	dcsClient := client.NewRedisClient()
	for _, instance := range viper.GetStringSlice("redis.instanceId") {
		req := &model.ListBackupRecordsRequest{}
		req.InstanceId = instance
		beginTimeRequest := viper.GetString("redis.beginTime")
		endTimeRequest := viper.GetString("redis.endTime")
		req.BeginTime = &beginTimeRequest
		req.EndTime = &endTimeRequest
		limitRequest := viper.GetInt32("redis.limit")
		req.Limit = &limitRequest
		response, err := dcsClient.ListBackupRecords(req)
		if err != nil {
			logger.Logger.Errorw("list backup error", "error", err.Error())
			// send notice to alarm center
			return
		}
		if response.HttpStatusCode == 200 {
			if response.BackupRecordResponse != nil {
				backupsList[instance] = *response.BackupRecordResponse
			}
		}
		//if len(*response.BackupRecordResponse) != 0 {
		//	backupsList[instance] = *response.BackupRecordResponse
		//}
		if len(backupsList) == 0 {
			logger.Logger.Infow("no manual backup found")
		}
	}
	for instance, backups := range backupsList {
		success := 0
		for _, backup := range backups {
			if backup.Status.Value() == "succeed" {
				request := &model.DeleteBackupFileRequest{}
				request.BackupId = *backup.BackupId
				request.InstanceId = instance
				_, err := dcsClient.DeleteBackupFile(request)
				if err != nil {
					logger.Logger.Errorw("delete redis backup error", "error", err.Error(), "instanceId", instance, "backupId", backup.BackupId, "beginTime", backup.CreatedAt)
					// send to alarm center
					continue
				}
				success++
				logger.Logger.Infow("[Success]delete redis backup success", "instanceId", instance, "backupName", backup.BackupName, "backupId", backup.BackupId, "beginTime", backup.CreatedAt)
			}

		}
		fmt.Printf("[Result]Total: %d,Success: %d\n", len(backups), success)
	}
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
	/* */
	//redisPusher := pusher.New(viper.GetString("global.pushgateway"), "huawei_redis_backup")
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
