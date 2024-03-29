package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChashFlow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Categories []Category         `bson:"Categories"`
}

type Transaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Name            string             `bson:"Name"`
	UserID          primitive.ObjectID `bson:"UserID"`
	TransactionType string             `bson:"TransactionType"`
	Category        string             `bson:"Category"`
	Duration        string             `bson:"Duration"`
	Amount          float32            `bson:"Amount"`
}

type Category struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"Name"`
	Transactions []Transaction      `bson:"Transactions"`
}
