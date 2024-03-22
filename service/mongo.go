package service

import (
	dds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dds/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dds/v3/model"
	"github.com/spf13/viper"
	"huawei_backup/logger"
	"huawei_backup/pkg/client"
	"huawei_backup/pkg/push"
	"sync"
)

type MongoBackup struct {
	Name string
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
