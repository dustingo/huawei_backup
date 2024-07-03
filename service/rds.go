package service

import (
	"fmt"
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
type RdsDeleteBackup struct {
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

	//	rdsPusher := pusher.New(viper.GetString("global.pushgateway"), "huawei_rds_backup")
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

func (r RdsDeleteBackup) Delete() {
	//var wg sync.WaitGroup
	var backupsList map[string][]model.BackupForList = make(map[string][]model.BackupForList)
	rdsClient := client.NewRdsClient()
	for _, instance := range viper.GetStringSlice("rds.instanceId") {
		req := &model.ListBackupsRequest{}
		req.InstanceId = instance
		backupTypeRequest := model.GetListBackupsRequestBackupTypeEnum().MANUAL
		req.BackupType = &backupTypeRequest
		limitRequest := viper.GetInt32("rds.limit")
		req.Limit = &limitRequest
		beginTimeRequest := viper.GetString("rds.beginTime")
		endTimeRequest := viper.GetString("rds.endTime")
		req.BeginTime = &beginTimeRequest
		req.EndTime = &endTimeRequest
		response, err := rdsClient.ListBackups(req)
		if err != nil {
			logger.Logger.Errorw("list backup error", "error", err.Error())
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
			_, err := rdsClient.DeleteManualBackup(request)
			if err != nil {
				logger.Logger.Errorw("delete manual backup error", "error", err.Error(), "instanceId", instance, "backupId", backup.Id, "beginTime", backup.BeginTime)
				// send to alarm center
				continue
			}
			success++
			logger.Logger.Infow("[Success]delete manual backup success", "instanceId", instance, "backupName", backup.Name, "backupId", backup.Id, "beginTime", backup.BeginTime)
		}
		fmt.Printf("[Result]Total: %d,Success: %d\n", len(backups), success)
	}

}
