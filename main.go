package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Document struct {
	Did     string    `bson:"did"`
	Utime   string    `bson:"utime"`
	Content []Content `bson:"content"`
}

type Content struct {
	Pid   string `bson:"pid"`
	Type  string `bson:"type"`
	Addr  string `bson:"addr"`
	Addrv string `bson:"addrv"`
	Ctime string `bson:"ctime"`
}

var (
	collection *mongo.Collection
	client     *mongo.Client
)

func NewFindOptions(pageNumber, pageSize int64) (option *options.FindOptions) {
	option = options.Find()
	option.SetLimit(int64(pageSize))
	option.SetSkip(int64((pageNumber - 1) * pageSize))
	return
}

func NewAggregateOptions(pageNumber, pageSize int64) (opts *options.AggregateOptions) {
	opts = options.Aggregate()
	return
}

func init() {

	var err error

	credential := options.Credential{
		Username: "root",
		Password: "123456",
	}

	url := "mongodb://82.157.202.19:27017"
	db := "plc_test"
	//collection := "report_data"

	clientOptions := options.Client().ApplyURI(url)

	// Set the Credential as the authentication option in the client options
	clientOptions.SetAuth(credential)

	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 选择数据库和集合
	database := client.Database(db)
	collection = database.Collection("report_data")

}

func InsertOne(ctx context.Context, collection *mongo.Collection, document interface{}, opts ...*options.InsertOneOptions) {
	collection.InsertOne(ctx, document, opts...)
}

func Find(ctx context.Context, collection *mongo.Collection, query interface{}, opts ...*options.FindOptions) {
	cursor, err := collection.Find(ctx, query, opts...)
	if err != nil {

	}
	var data []Document
	cursor.All(ctx, &data)
	fmt.Println("len data ==== ", len(data))
	fmt.Println("data ==== ", data)
}

func Aggregate(ctx context.Context, collection *mongo.Collection, pipeline interface{},
	data interface{}, opts ...*options.AggregateOptions) (err error) {
	cursor, err := collection.Aggregate(context.Background(), pipeline, opts...)
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(ctx, data)
	fmt.Println("======", data)
	return err
}

func Update(ctx context.Context, collection *mongo.Collection, query interface{},
	data interface{}) (result *mongo.UpdateResult, err error) {
	result, err = collection.UpdateOne(ctx, query, data)
	return
}

func main() {
	ctx := context.Background()

	// 构建聚合管道
	/*
		pipeline := mongo.Pipeline{
			{{"$match", bson.M{"content.pid": "1"}}},
			{{"$project", bson.M{
				"did":   1,
				"utime": 1,
				"content": bson.M{
					"$filter": bson.M{
						"input": "$content",
						"as":    "c",
						"cond":  bson.M{"$eq": []interface{}{"$$c.pid", "1"}},
					},
				},
			}}},
		}
	*/

	pageSize := 10  // 每页数据数量
	pageNumber := 1 // 页码

	// 计算跳过的文档数量
	skip := (pageNumber - 1) * pageSize

	// 构建聚合管道
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"did":         "FHC230980939",
			"content.pid": "1",
		}}},
		{{"$project", bson.M{
			"did":   1,
			"utime": 1,
			"content": bson.M{
				"$filter": bson.M{
					"input": "$content",
					"as":    "c",
					"cond":  bson.M{"$eq": []interface{}{"$$c.pid", "1"}},
				},
			},
		}}},
		{{"$skip", skip}},
		{{"$limit", pageSize}},
	}

	var data []Document

	Aggregate(ctx, collection, pipeline, &data)

	fmt.Println(data)

}
