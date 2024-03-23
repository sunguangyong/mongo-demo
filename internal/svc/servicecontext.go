package svc

import (
	"mongo-demo/internal/config"
	"mongo-demo/model"
)

type ServiceContext struct {
	MongoReportData model.MongoReportDataModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		MongoReportData: model.NewMongoReportDataModel(c),
	}
}
