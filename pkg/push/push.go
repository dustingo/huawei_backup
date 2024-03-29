package push

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/viper"
	"huawei_backup/logger"
)

func Push(db, instance string) {
	if !viper.GetBool("push") {
		logger.Logger.Infow("not push to prometheus pushgateway", "db", db, "instance", instance)
		return
	}
	var p prometheus.Gauge
	var pusher *push.Pusher

	switch db {
	case "redis":
		p = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "redis_backup_execution",
			Help: "redis backup manual execution",
		})
		p.Set(1)
		pusher = push.New(viper.GetString("global.pushgateway"), "huawei_redis_backup").Collector(p).Grouping("instanceId", instance).Grouping("db", db)
	case "mongo":
		p = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "mongo_backup_execution",
			Help: "mongo backup manual execution",
		})
		//p.WithLabelValues(db, instance).Set(1)
		p.Set(1)
		pusher = push.New(viper.GetString("global.pushgateway"), "huawei_mongo_backup").Collector(p).Grouping("instanceId", instance).Grouping("db", db)
	case "rds":
		p = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "rds_backup_execution",
			Help: "rds backup manual execution",
		})
		//p.WithLabelValues(db, instance).Set(1)
		p.Set(1)
		pusher = push.New(viper.GetString("global.pushgateway"), "huawei_rds_backup").Collector(p).Grouping("instanceId", instance).Grouping("db", db)
	}
	//pusher.Collector(p)
	if err := pusher.Push(); err != nil {
		logger.Logger.Errorw("push to prometheus pushgateway error", "error", err.Error())
		return
	}
	logger.Logger.Infow("push to prometheus pushgateway success", "db", db, "instanceId", instance)
}
