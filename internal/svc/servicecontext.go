package svc

import (
	"mongo-demo/internal/config"
	"mongo-demo/model/mongo"
)

type ServiceContext struct {
	MongoReportDataConn mongo.MongoReportDataModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		MongoReportDataConn: mongo.NewMongoReportDataModel(c.ReportDataMongo),
	}
}
