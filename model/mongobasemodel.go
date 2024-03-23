package model

import (
	"context"
	"log"
	"mongo-demo/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCollection(mongoConf config.MongoConf) (collection *mongo.Collection) {
	var err error

	clientOptions := options.Client().ApplyURI(mongoConf.Url)
	clientOptions.SetAuth(mongoConf.Credential)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// 选择数据库和集合
	database := client.Database(mongoConf.DbName)
	collection = database.Collection(mongoConf.CollectionName)
	return
}

func NewFindOptions(pageNumber, pageSize int64) (option *options.FindOptions) {
	option = options.Find()
	option.SetLimit(int64(pageSize))
	option.SetSkip(int64((pageNumber - 1) * pageSize))
	return
}
