package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"finance-app/models"
)

// Fungsi helper untuk koneksi ke MongoDB
func connectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping server untuk memastikan koneksi berhasil
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Controllers untuk Kategori

func createCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	_, err = collection.InsertOne(context.TODO(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func getCategories(c *gin.Context) {
	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var categories []models.Category
	for cursor.Next(context.TODO()) {
		var category models.Category
		err := cursor.Decode(&category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func updateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var update bson.M
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	filter := bson.M{"_id": id}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func deleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("categories")

	filter := bson.M{"_id": id}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// Controllers untuk Transaksi

func createTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("transactions")

	_, err = collection.InsertOne(context.TODO(), transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func getTransactions(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	transactions, err := getTransactionsByDateRange(startDate, endDate, client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func getTransactionsByDateRange(startDate, endDate time.Time, client *mongo.Client) ([]models.Transaction, error) {
	collection := client.Database("finance_app").Collection("transactions")

	// Jika startDate dan endDate kosong, gunakan tanggal awal dan akhir bulan ini
	if startDate.IsZero() && endDate.IsZero() {
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		startDate = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local) // Removed the extra 0
		endDate = startDate.AddDate(0, 1, -1)
	}

	// Buat filter query
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	// Eksekusi query
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterasi hasil query
	var transactions []models.Transaction
	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		err := cursor.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func updateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var update bson.M
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("transactions")

	filter := bson.M{"_id": id}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

func deleteTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("finance_app").Collection("transactions")

	filter := bson.M{"_id": id}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// Fungsi untuk mendapatkan saldo saat ini
func getCurrentBalance(c *gin.Context) {
	client, err := connectToMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	balance, err := calculateCurrentBalance(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func calculateCurrentBalance(client *mongo.Client) (float64, error) {
	collection := client.Database("finance_app").Collection("transactions")

	// Aggregate untuk menghitung saldo
	pipeline := mongo.Pipeline{
		bson.D{{"$group", bson.D{
			{"_id", nil},
			{"balance", bson.D{{"$sum", bson.D{{"$cond", bson.A{
				bson.M{"$eq": bson.A{"$type", "income"}},
				"$amount",
				bson.M{"$multiply": bson.A{"$amount", -1}},
			}}}}},
			}}},
		},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil // Tidak ada transaksi, saldo 0
	}

	return result[0]["balance"].(float64), nil
}

func main() {
	r := gin.Default()

	// Endpoint untuk Kategori
	r.POST("/categories", createCategory)
	r.GET("/categories", getCategories)
	r.PUT("/categories/:id", updateCategory)
	r.DELETE("/categories/:id", deleteCategory)

	// Endpoint untuk Transaksi
	r.POST("/transactions", createTransaction)
	r.GET("/transactions", getTransactions)
	r.PUT("/transactions/:id", updateTransaction)
	r.DELETE("/transactions/:id", deleteTransaction)

	// Endpoint untuk Saldo
	r.GET("/balance", getCurrentBalance)

	r.Run()
}
