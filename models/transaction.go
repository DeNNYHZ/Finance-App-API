package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Type        string             `bson:"type"` // "income" atau "expense"
	CategoryID  primitive.ObjectID `bson:"category_id"`
	Amount      float64            `bson:"amount"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
	UserID      primitive.ObjectID `bson:"user_id"`
}
