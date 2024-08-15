package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Type        string             `bson:"type"`    // "income" atau "expense"
	UserID      primitive.ObjectID `bson:"user_id"` // Tambahkan field ini jika diperlukan
}
