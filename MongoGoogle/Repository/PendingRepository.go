package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePendingInvite(email string, companyId string) (model.PendingRequest, primitive.ObjectID) {
	UsersCollection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	identifier, iderr := primitive.ObjectIDFromHex(companyId)
	if iderr != nil {
		log.Fatal(iderr)
	}
	// Creating user instance
	request := model.PendingRequest{
		Email:     email,
		CompanyID: identifier,
		Completed: false,
	}
	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added new request with ID:", insertResult.InsertedID)
	id := insertResult.InsertedID.(primitive.ObjectID)
	return request, id
}
