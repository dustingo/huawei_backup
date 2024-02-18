package service

import (
	rdsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	"github.com/spf13/viper"
	"huawei_backup/logger"
	"huawei_backup/pkg/client"
	"huawei_backup/pkg/push"
	"sync"
)

type RdsBackup struct {
	Name string
}

func (r RdsBackup) Backup() {
	var wg sync.WaitGroup
	rdsClient := client.NewRdsClient()
	requests := make([]*model.CreateManualBackupRequest, 0)
	for _, instance := range viper.GetStringSlice("rds.instanceId") {
		req := &model.CreateManualBackupRequest{}
		body := &model.CreateManualBackupRequestBody{
			Name:       r.Name,
			InstanceId: instance,
		}
		req.Body = body
		requests = append(requests, req)
	}
	for _, request := range requests {
		wg.Add(1)
		go func(client *rdsv3.RdsClient, req *model.CreateManualBackupRequest) {
			defer wg.Done()
			response, err := client.CreateManualBackup(req)
			if err != nil {
				logger.Logger.Errorw("create manual backup error", "error", err.Error())
				return
			}
			push.Push("rds", response.Backup.InstanceId)
			logger.Logger.Infow("create manual backup success", "InstanceId", response.Backup.InstanceId, "BeginTime", response.Backup.BeginTime, "Type", response.Backup.Type)
		}(rdsClient, request)
	}
	wg.Wait()
}
