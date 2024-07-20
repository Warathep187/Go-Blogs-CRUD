package connections

import (
	"context"
	"fmt"
	"go_blogs/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newMongoClient() *mongo.Client {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s",
		configs.Env.MongoUsername,
		configs.Env.MongoPassword,
		configs.Env.MongoEndpoint,
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Mongo connected")
	return client
}

func NewMongoCollection(database *mongo.Database, collectionName string) *mongo.Collection {
	ctx := context.TODO()
	err := database.CreateCollection(ctx, collectionName)
	if err != nil {
		panic(err)
	}

	collection := database.Collection(collectionName)
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
	})
	if err != nil {
		panic(err)
	}
	return collection
}
