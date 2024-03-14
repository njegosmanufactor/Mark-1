package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Country string `bson:"Country"`
	City    string `bson:"City"`
	Street  string `bson:"Street"`
}

type Company struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	Name                  string             `bson:"Name"`
	Address               Location           `bson:"Address"`
	Website               string             `bson:"Website"`
	ListOfApprovedDomains []string           `bson:"Domains"`
	Owner                 primitive.ObjectID `bson:"ownerId"`
	Employees             []string           `bson:"Employees"`
}
