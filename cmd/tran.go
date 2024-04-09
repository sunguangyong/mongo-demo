package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const collectionName = "your_collection"

// 事务事例
func main() {
	// 连接到MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("your_mongodb_uri"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// 获取集合
	collection := client.Database("your_database").Collection(collectionName)

	// 开始事务
	session, err := client.StartSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.EndSession(ctx)

	// 开始事务
	err = session.StartTransaction()
	if err != nil {
		log.Fatal(err)
	}

	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {

		doc1 := bson.D{{"name", "example"}, {"age", 25}}
		_, err = collection.InsertOne(ctx, doc1)

		doc2 := bson.D{{"name", "example"}, {"age", 25}}
		_, err = collection.InsertOne(ctx, doc2)

		return err

	})
	if err != nil {
		session.AbortTransaction(ctx)
	}

	// 提交事务
	err = session.CommitTransaction(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Document inserted successfully within a transaction")
}
