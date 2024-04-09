package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mongo-demo/internal/config"
	"mongo-demo/internal/svc"
	"mongo-demo/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var svctx *svc.ServiceContext

func init() {
	var mongoConf config.MongoConf
	mongoConf.Url = "mongodb://82.157.202.19:27017"
	mongoConf.DbName = "plc_test"
	mongoConf.CollectionName = "report_data"
	mongoConf.Credential = options.Credential{
		Username: "root",
		Password: "123456",
	}

	c := config.Config{
		ReportDataMongo: mongoConf,
	}
	svctx = svc.NewServiceContext(c)
}

func main() {
	ctx := context.Background()
	AggregateDemo(ctx)
	FindDemo(ctx)
	//InsertDemo(ctx)
	//UpdateDemo(ctx)
}

func FindDemo(ctx context.Context) {
	var query interface{}

	optionLimit := mongo.NewFindOptions(1, 2)
	optionOrder := mongo.NewSortOptions("utime", 1)

	query = bson.M{}

	data, err := svctx.MongoReportDataConn.Find(ctx, query, optionLimit, optionOrder)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
}

func AggregateDemo(ctx context.Context) {
	pipeline := svctx.MongoReportDataConn.NewPidPipLine(ctx, "1", "FHC230980939", 1, 1)

	data, err := svctx.MongoReportDataConn.Aggregate(ctx, pipeline)

	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(data)
}

func InsertDemo(ctx context.Context) {
	data := `{"did":"FH5260674584","utime":"2023/10/11 15:06:50","content":[{"pid":"04","type":"0","addr":"电流","addrv":"4614","ctime":"2023/10/11 15:06:50"}]}`

	var document mongo.ReportDataModel

	json.Unmarshal([]byte(data), &document)

	svctx.MongoReportDataConn.InsertOne(ctx, &document)
}

func UpdateDemo(ctx context.Context) {
	// 更新条件
	filter := bson.M{"did": "456"}

	data, err := svctx.MongoReportDataConn.Find(ctx, filter)

	if err != nil {
		fmt.Println("err === ", err)
	}

	if len(data) > 0 {
		doc := data[0]
		doc.Content[0].Pid = "100"
		result, err := svctx.MongoReportDataConn.Update(ctx, filter, doc)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}
}
