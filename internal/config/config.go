package config

import "go.mongodb.org/mongo-driver/mongo/options"

type MongoConf struct {
	Credential     options.Credential `json:"credential"`
	Url            string             `json:"url"`
	DbName         string             `json:"db_name"`
	CollectionName string             `json:"collection_name"`
}

type Config struct {
	ReportDataMongo MongoConf
}
