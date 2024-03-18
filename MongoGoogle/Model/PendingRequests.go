package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PendingRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"Email,omitempty"`
	CompanyID primitive.ObjectID `bson:"companyId,omitempty"`
	Completed bool               `bson:"Completed"`
}
type PendingOwnershipTransfer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"Email,omitempty"`
	CompanyID primitive.ObjectID `bson:"CompanyId,omitempty"`
	OwnerID   primitive.ObjectID `bson:"OwnerId,omitempty"`
	Completed bool               `bson:"Completed"`
}

type PasswordChangeRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"Email,omitempty"`
	Completed bool               `bson:"Completed"`
}
