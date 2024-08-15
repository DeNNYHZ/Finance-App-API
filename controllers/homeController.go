package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"finance-app/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetHomeData(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(primitive.ObjectID)

	// 1. Hitung Saldo Saat Ini
	var currentBalance float64
	matchStage := bson.D{{"$match", bson.D{{"user_id", userID}}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id", nil},
		{"totalIncome", bson.D{{"$sum", bson.D{{"$cond", bson.A{
			bson.D{{"$eq", bson.A{"$type", "income"}}},
			"$amount",
			0,
		}}}}}},
		{"totalExpense", bson.D{{"$sum", bson.D{{"$cond", bson.A{
			bson.D{{"$eq", bson.A{"$type", "expense"}}},
			"$amount",
			0,
		}}}}}},
	}}}

	cursor, err := database.TransactionCollection.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		http.Error(w, "Error calculating current balance", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		http.Error(w, "Error decoding balance result", http.StatusInternalServerError)
		return
	}

	if len(result) > 0 {
		currentBalance = result[0]["totalIncome"].(float64) - result[0]["totalExpense"].(float64)
	}

	// 2. Hitung Total Pengeluaran (All Time)
	filter := bson.M{"user_id": userID, "type": "expense"}
	totalExpense, err := calculateTotalAmount(filter)
	if err != nil {
		http.Error(w, "Error calculating total expense", http.StatusInternalServerError)
		return
	}

	// 3. Hitung Total Pemasukan (All Time)
	filter = bson.M{"user_id": userID, "type": "income"}
	totalIncome, err := calculateTotalAmount(filter)
	if err != nil {
		http.Error(w, "Error calculating total income", http.StatusInternalServerError)
		return
	}

	// Kirim respons
	homeData := map[string]interface{}{
		"current_balance": currentBalance,
		"total_expense":   totalExpense,
		"total_income":    totalIncome,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(homeData)
}

// Fungsi helper untuk menghitung total amount berdasarkan filter
func calculateTotalAmount(filter bson.M) (float64, error) {
	groupStage := bson.D{{"$group", bson.D{
		{"_id", nil},
		{"total", bson.D{{"$sum", "$amount"}}},
	}}}

	cursor, err := database.TransactionCollection.Aggregate(context.Background(), mongo.Pipeline{
		bson.D{{"$match", filter}},
		groupStage,
	})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.Background())

	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}

	if len(result) > 0 {
		return result[0]["total"].(float64), nil
	}
	return 0, nil
}
