package model

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"log"
)

type MongoConf struct {
	Credential     options.Credential `json:"credential"`
	Url            string             `json:"url"`
	DbName         string             `json:"db_name"`
	CollectionName string             `json:"collection_name"`
}

type defaultMongoReportDataModel struct {
	Collection *mongo.Collection
}

type MongoReportDataModel interface {
	InsertOne(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions)

	Find(ctx context.Context, data interface{}, query interface{},
		opts ...*options.FindOptions)

	Aggregate(ctx context.Context, pipeline interface{},
		data interface{}, opts ...*options.AggregateOptions) (err error)

	Update(ctx context.Context, query interface{},
		data interface{}) (result *mongo.UpdateResult, err error)
}

func NewMongoReportDataModel(mongoConf MongoConf) MongoReportDataModel {
	return &defaultMongoReportDataModel{
		Collection: NewCollection(mongoConf),
	}
}

func (m *defaultMongoReportDataModel) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) {
	m.Collection.InsertOne(ctx, document, opts...)
}

func (m *defaultMongoReportDataModel) Find(ctx context.Context, data interface{}, query interface{},
	opts ...*options.FindOptions) {
	cursor, err := m.Collection.Find(ctx, query, opts...)
	if err != nil {

	}
	cursor.All(ctx, data)
	fmt.Println("data ==== ", data)
}

func (m *defaultMongoReportDataModel) Aggregate(ctx context.Context, pipeline interface{},
	data interface{}, opts ...*options.AggregateOptions) (err error) {
	cursor, err := m.Collection.Aggregate(context.Background(), pipeline, opts...)
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(ctx, data)
	fmt.Println("======", data)
	return err
}

func (m *defaultMongoReportDataModel) Update(ctx context.Context, query interface{},
	data interface{}) (result *mongo.UpdateResult, err error) {
	result, err = m.Collection.UpdateOne(ctx, query, data)
	return
}
