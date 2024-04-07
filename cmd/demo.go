package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomHook struct {
}

func (h *CustomHook) CommandFinishedEvent(ctx context.Context, cevent *event.CommandFinishedEvent) {
	fmt.Println("ConnectionCreated event:", cevent)
	// 在连接创建时执行自定义逻辑
}

func (h *CustomHook) CommandFailedEvent(ctx context.Context, cevent *event.CommandFailedEvent) {
	fmt.Println("ConnectionCreated event:", cevent)
	// 在连接创建时执行自定义逻辑
}

func (h *CustomHook) CommandSucceededEvent(ctx context.Context, cevent *event.CommandSucceededEvent) {
	fmt.Println("ConnectionClosed event:", cevent)
	// 命令执行成功时自定义逻辑
}

func (h *CustomHook) CommandStarted(ctx context.Context, cevent *event.CommandStartedEvent) {
	fmt.Println("CommandStarted event:", cevent)
	// 在命令开始时执行自定义逻辑
}

func main() {
	// 连接 MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// 注册自定义钩子
	hook := &event.CommandMonitor{}
	clientOpts := options.Client().SetMonitor(hook)

	// 获取数据库和集合
	db := client.Database("mydb")
	collection := db.Collection("mycollection")

	// 执行数据库操作
	_, err = collection.InsertOne(context.Background(), map[string]interface{}{
		"field": "value",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 断开与 MongoDB 的连接
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
