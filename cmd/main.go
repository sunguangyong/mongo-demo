package cmd

import (
	"mongo-demo/model"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var mongoConf model.MongoConf
	mongoConf.Url = "mongodb://82.157.202.19:27017"
	mongoConf.DbName = "plc_test"
	mongoConf.CollectionName = "report_data"
	mongoConf.Credential = options.Credential{
		Username: "root",
		Password: "123456",
	}

	collection := model.NewCollection(mongoConf)

}
