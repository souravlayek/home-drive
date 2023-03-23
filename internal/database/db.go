package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MetaDataCollection *mongo.Collection

func ConnectDB() {
	dbName := "storage"
	mongoURL := os.Getenv("MONGO_URL")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongoURL).
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
