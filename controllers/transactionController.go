package controllers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"

	"finance-app/database"
	"finance-app/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari context
	userID := r.Context().Value("user_id").(primitive.ObjectID)

	// Dapatkan parameter filter dari query string (jika ada)
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Buat filter berdasarkan tanggal (jika ada parameter)
	filter := bson.M{"user_id": userID}
	if startDateStr != "" && endDateStr != "" {
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		endDate, _ := time.Parse("2006-01-02", endDateStr)
		filter["date"] = bson.M{"$gte": startDate, "$lte": endDate}
	} else {
		// Jika tidak ada filter tanggal, tampilkan transaksi bulan ini
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		filter["date"] = bson.M{"$gte": firstOfMonth, "$lte": lastOfMonth}
	}

	// Query semua transaksi milik user dengan filter
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", -1}}) // Urutkan berdasarkan tanggal terbaru

	cursor, err := database.TransactionCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var transactions []models.Transaction
	if err = cursor.All(context.Background(), &transactions); err != nil {
		http.Error(w, "Error decoding transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari context
	userID := r.Context().Value("user_id").(primitive.ObjectID)

	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transaction.UserID = userID // Set user ID pada transaksi
	transaction.Date = time.Now()

	// Validasi tipe transaksi
	if transaction.Type != "income" && transaction.Type != "expense" {
		http.Error(w, "Invalid transaction type. Must be 'income' or 'expense'", http.StatusBadRequest)
		return
	}

	// Simpan transaksi ke database
	result, err := database.TransactionCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		http.Error(w, "Error creating transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"inserted_id": result.InsertedID})
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari context
	userID := r.Context().Value("user_id").(primitive.ObjectID)

	// Ambil transactionID dari URL params
	params := mux.Vars(r)
	transactionID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	// Validasi tipe transaksi
	if transaction.Type != "income" && transaction.Type != "expense" {
		http.Error(w, "Invalid transaction type. Must be 'income' or 'expense'", http.StatusBadRequest)
		return
	}

	// Update transaksi di database (pastikan hanya transaksi milik user yang diupdate)
	filter := bson.M{"_id": transactionID, "user_id": userID}
	update := bson.M{"$set": bson.M{
		"type":        transaction.Type,
		"category_id": transaction.CategoryID,
		"amount":      transaction.Amount,
		"description": transaction.Description,
	}}
	result, err := database.TransactionCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating transaction", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Transaction not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction updated successfully"})
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	// Ambil user ID dari context
	userID := r.Context().Value("user_id").(primitive.ObjectID)

	params := mux.Vars(r)
	transactionID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	// Hapus transaksi dari database (pastikan hanya transaksi milik user yang dihapus)
	filter := bson.M{"_id": transactionID, "user_id": userID}
	result, err := database.TransactionCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error deleting transaction", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Transaction not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted successfully"})
}
