package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MetaDataCollection *mongo.Collection

func ConnectDB() {
	dbName := "storage"

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	MetaDataCollection = client.Database(dbName).Collection("MetaData")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected successfully")

}
