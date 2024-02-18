package client

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	dcs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dcs/v2"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dcs/v2/region"
	dds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dds/v3"
	ddsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dds/v3/region"
	rdsv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	regionv3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/region"
	"github.com/spf13/viper"
)

func NewRdsClient() *rdsv3.RdsClient {

	auth := basic.NewCredentialsBuilder().
		WithAk(viper.GetString("global.ak")).
		WithSk(viper.GetString("global.sk")).
		WithProjectId(viper.GetString("global.projectId")).
		Build()
	return rdsv3.NewRdsClient(
		rdsv3.RdsClientBuilder().
			WithRegion(regionv3.ValueOf(viper.GetString("global.region"))).
			WithCredential(auth).
			Build())
}
func NewRedisClient() *dcs.DcsClient {
	auth := basic.NewCredentialsBuilder().
		WithAk(viper.GetString("global.ak")).
		WithSk(viper.GetString("global.sk")).
		WithProjectId(viper.GetString("global.projectId")).
		Build()
	return dcs.NewDcsClient(
		dcs.DcsClientBuilder().
			WithRegion(region.ValueOf(viper.GetString("global.region"))).
			WithCredential(auth).
			Build())

}

func NewMongoClient() *dds.DdsClient {
	auth := basic.NewCredentialsBuilder().
		WithAk(viper.GetString("global.ak")).
		WithSk(viper.GetString("global.sk")).
		WithProjectId(viper.GetString("global.projectId")).
		Build()
	return dds.NewDdsClient(
		dds.DdsClientBuilder().
			WithRegion(ddsRegion.ValueOf(viper.GetString("global.region"))).
			WithCredential(auth).
			Build())
}
