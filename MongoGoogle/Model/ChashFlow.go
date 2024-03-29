package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CashFlow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Categories []Category         `bson:"Categories"`
}
type Category struct {
	Name          string        `bson:"Name"`
	Subcategories []Subcategory `bson:"Subcategory"`
}
type Subcategory struct {
	Name     string        `bson:"Name"`
	Positive []Transaction `bson:"Inflow"`
	Negative []Transaction `bson:"Outflow"`
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
