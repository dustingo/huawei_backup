package service

import (
	"fmt"
	dds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dds/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dds/v3/model"
	"github.com/spf13/viper"
	"huawei_backup/logger"
	"huawei_backup/pkg/client"
	"huawei_backup/pkg/push"
	"sync"
	"time"
)

type MongoBackup struct {
	Name string
}
type MongoDeleteBackup struct {
	Name string
}

func (m MongoDeleteBackup) Delete() {
	var backupsList map[string][]model.BackupForList = make(map[string][]model.BackupForList)
	mClient := client.NewMongoClient()
	for _, instance := range viper.GetStringSlice("mongo.instanceId") {
		req := &model.ListBackupsRequest{}
		req.InstanceId = &instance
		backupTypeRequest := model.GetListBackupsRequestBackupTypeEnum().MANUAL
		req.BackupType = &backupTypeRequest
		limitRequest := viper.GetInt32("mongo.limit")
		req.Limit = &limitRequest
		beginTimeRequest := viper.GetString("mongo.beginTime")
		endTimeRequest := viper.GetString("mongo.endTime")
		req.BeginTime = &beginTimeRequest
		req.EndTime = &endTimeRequest
		response, err := mClient.ListBackups(req)
		if err != nil {
			logger.Logger.Errorw("list mongodb backup error", "instanceId", instance, "error", err.Error())
			// send notice to alarm center
			return
		}
		if len(*response.Backups) != 0 {
			backupsList[instance] = *response.Backups
		}
	}
	if len(backupsList) == 0 {
		logger.Logger.Infow("no manual backup found")
	}
	for instance, backups := range backupsList {
		success := 0
		for _, backup := range backups {
			request := &model.DeleteManualBackupRequest{}
			request.BackupId = backup.Id
			_, err := mClient.DeleteManualBackup(request)
			if err != nil {
				logger.Logger.Errorw("delete mongo manual backup error", "error", err.Error(), "instanceId", instance, "backupId", backup.Id, "beginTime", backup.BeginTime)
				// send to alarm center
				continue
			}
			success++
			logger.Logger.Infow("[Success]delete mongo manual backup success", "instanceId", instance, "instanceName", backup.InstanceName, "backupId", backup.Id, "beginTime", backup.BeginTime)
			time.Sleep(time.Second * 1)
		}
		fmt.Printf("[Result]Total: %d,Success: %d\n", len(backups), success)
	}

}

func (m MongoBackup) Backup() {
	var wg sync.WaitGroup
	mongoClient := client.NewMongoClient()
	requests := make([]*model.CreateManualBackupRequest, 0)
	for _, instance := range viper.GetStringSlice("mongo.instanceId") {
		req := &model.CreateManualBackupRequest{}
		option := &model.CreateManualBackupOption{
			Name:       m.Name,
			InstanceId: instance,
		}
		req.Body = &model.CreateManualBackupRequestBody{
			Backup: option,
		}
		requests = append(requests, req)
	}
	//mongoPusher := pusher.New(viper.GetString("global.pushgateway"), "huawei_mongo_backup")
	for _, request := range requests {
		wg.Add(1)
		go func(client *dds.DdsClient, req *model.CreateManualBackupRequest) {
			defer wg.Done()
			response, err := client.CreateManualBackup(req)
			if err != nil {
				logger.Logger.Errorw("create mongodb manual backup error", "error", err.Error())
				return
			}
			push.Push("mongo", req.Body.Backup.InstanceId)
			logger.Logger.Infow("create mongodb manual backup success", "InstanceId", req.Body.Backup.InstanceId, "BackupId", response.BackupId)
		}(mongoClient, request)
	}
	wg.Wait()
}
