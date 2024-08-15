package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"finance-app/models"
)

func createCategory(category models.Category) error {
	// Koneksi ke MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	_, err = collection.InsertOne(context.TODO(), category)
	return err
}

func getCategories() ([]models.Category, error) {
	// Koneksi ke MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var categories []models.Category
	for cursor.Next(context.TODO()) {
		var category models.Category
		err := cursor.Decode(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func updateCategory(id primitive.ObjectID, update bson.M) error {
	// Koneksi ke MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	filter := bson.M{"_id": id}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": update})
	return err
}

func deleteCategory(id primitive.ObjectID) error {
	// Koneksi ke MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	filter := bson.M{"_id": id}

	_, err = collection.DeleteOne(context.TODO(), filter)
	return err
}
