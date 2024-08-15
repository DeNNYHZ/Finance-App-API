package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client                *mongo.Client
	UserCollection        *mongo.Collection
	CategoryCollection    *mongo.Collection
	TransactionCollection *mongo.Collection
)

func ConnectDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Ganti dengan URI MongoDB Anda
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping server untuk memastikan koneksi berhasil
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Pilih database
	db := client.Database("finance_app")

	// Dapatkan koleksi (collections)
	UserCollection = db.Collection("users")
	CategoryCollection = db.Collection("categories")
	TransactionCollection = db.Collection("transactions")

	log.Println("Connected to MongoDB!")
	return client, nil
}
