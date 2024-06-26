package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CashFlow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Categories []Category         `bson:"Categories"`
}
type Category struct {
	Name          string        `bson:"Name"`
	Subcategories []Subcategory `bson:"Subcategories"`
}
type Subcategory struct {
	Name    string        `bson:"Name"`
	Inflow  []Transaction `bson:"Inflow"`
	Outflow []Transaction `bson:"Outflow"`
}
type Transaction struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"Name"`
	UserID   primitive.ObjectID `bson:"UserID"`
	Duration string             `bson:"Duration"`
	Amount   float32            `bson:"Amount"`
}
