package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://root:example@localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	database := client.Database("test")
	collection := database.Collection("test")

	// Insert
	res, err := collection.InsertOne(ctx, map[string]string{"name": "pi", "value": "3.14159"})
	if err != nil {
		panic(err)
	}

	id := res.InsertedID.(primitive.ObjectID)
	println(id.Hex())
}
