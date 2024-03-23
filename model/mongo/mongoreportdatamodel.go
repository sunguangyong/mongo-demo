package mongo

import (
	"encoding/json"
	"fmt"
	"mongo-demo/internal/config"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
)

type defaultMongoReportDataModel struct {
	Collection *mongo.Collection
}

type ReportDataModel struct {
	Did     string `json:"did"`
	Utime   string `json:"utime"`
	Content []struct {
		Pid   string `json:"pid"`
		Type  string `json:"type"`
		Addr  string `json:"addr"`
		Addrv string `json:"addrv"`
		Ctime string `json:"ctime"`
	} `json:"content"`
}

type MongoReportDataModel interface {
	InsertOne(ctx context.Context, document *ReportDataModel,
		opts ...*options.InsertOneOptions)

	Find(ctx context.Context, query interface{},
		opts ...*options.FindOptions) (data []*ReportDataModel, err error)

	Aggregate(ctx context.Context, pipeline interface{},
		opts ...*options.AggregateOptions) (data []*ReportDataModel, err error)

	NewPidPipLine(ctx context.Context, pid, did string, pageNumber,
		pageSize int64) (pipeline mongo.Pipeline)

	Update(ctx context.Context, query interface{},
		data *ReportDataModel) (result *mongo.UpdateResult, err error)
}

func NewMongoReportDataModel(mongoConf config.MongoConf) MongoReportDataModel {
	return &defaultMongoReportDataModel{
		Collection: NewCollection(mongoConf),
	}
}

func (m *defaultMongoReportDataModel) InsertOne(ctx context.Context, document *ReportDataModel,
	opts ...*options.InsertOneOptions) {
	m.Collection.InsertOne(ctx, document, opts...)
}

func (m *defaultMongoReportDataModel) Find(ctx context.Context, query interface{},
	opts ...*options.FindOptions) (data []*ReportDataModel, err error) {
	cursor, err := m.Collection.Find(ctx, query, opts...)

	if err != nil {
		return
	}
	cursor.All(ctx, &data)
	return
}

func (m *defaultMongoReportDataModel) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (data []*ReportDataModel, err error) {
	cursor, err := m.Collection.Aggregate(context.Background(), pipeline, opts...)
	if err != nil {
		return
	}
	err = cursor.All(ctx, &data)
	return
}

func (m *defaultMongoReportDataModel) NewPidPipLine(ctx context.Context, pid, did string, pageNumber,
	pageSize int64) (pipeline mongo.Pipeline) {

	pipeline = mongo.Pipeline{
		{{"$match", bson.M{
			"did":         did,
			"content.pid": pid,
		}}},
		{{"$project", bson.M{
			"did":   1,
			"utime": 1,
			"content": bson.M{
				"$filter": bson.M{
					"input": "$content",
					"as":    "c",
					"cond":  bson.M{"$eq": []interface{}{"$$c.pid", pid}},
				},
			},
		}}},
		{{"$sort", bson.M{"utime": -1}}},
		//{{"$skip", skip}},
		//{{"$limit", pageSize}},
	}

	if pageNumber > 0 && pageSize > 0 {
		// 计算跳过的文档数量
		skip := (pageNumber - 1) * pageSize
		pipeline = append(pipeline, bson.D{
			{"$skip", skip},
		})

		pipeline = append(pipeline, bson.D{
			{"$limit", pageSize},
		})

		//{"$skip", skip},
		//{"$limit", pageSize},
	}
	return pipeline
}

func (m *defaultMongoReportDataModel) Update(ctx context.Context, query interface{},
	data *ReportDataModel) (result *mongo.UpdateResult, err error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal report:", err)
		return
	}

	var bsonData bson.M
	err = bson.UnmarshalExtJSON(jsonData, true, &bsonData)

	a, err := bson.Marshal(data)
	fmt.Println("aaaaaa", string(a))
	fmt.Println("aaaaaa", err)

	result, err = m.Collection.UpdateMany(ctx, query, bson.M{"$set": bsonData})
	return
}

func (m *defaultMongoReportDataModel) GetUpdateData(data *ReportDataModel) (updateData bson.M) {

	//updateData = bson.M{
	//	"$set": bson.M{"did": data.Did,"utime": data.Utime, },
	//	"$set": bson.M{},
	//	"$set": bson.M{"content": data.Content},
	//	}
	return
}
