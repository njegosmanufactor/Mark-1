package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePasswordChangeRequest(email string) (model.PasswordChangeRequest, primitive.ObjectID) {
	RequestCollection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	request := model.PasswordChangeRequest{
		Email:     email,
		Completed: false,
	}
	insertResult, err := RequestCollection.InsertOne(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added new request with ID:", insertResult.InsertedID)
	id := insertResult.InsertedID.(primitive.ObjectID)
	return request, id
}

func CreatePendingInvite(email string, companyId string) (model.PendingRequest, primitive.ObjectID) {
	RequestCollection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	identifier, iderr := primitive.ObjectIDFromHex(companyId)
	if iderr != nil {
		log.Fatal(iderr)
	}
	// Creating request instance
	request := model.PendingRequest{
		Email:     email,
		CompanyID: identifier,
		Completed: false,
	}
	// Adding request to the database
	insertResult, err := RequestCollection.InsertOne(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added new request with ID:", insertResult.InsertedID)
	id := insertResult.InsertedID.(primitive.ObjectID)
	return request, id
}

func CreatePendingOwnershipInvitation(email string, ownerId primitive.ObjectID, companyId primitive.ObjectID) (model.PendingOwnershipTransfer, primitive.ObjectID) {
	RequestCollection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	// Creating request instance
	request := model.PendingOwnershipTransfer{
		Email:     email,
		CompanyID: companyId,
		OwnerID:   ownerId,
		Completed: false,
	}
	// Adding request to the database
	insertResult, err := RequestCollection.InsertOne(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added new request with ID:", insertResult.InsertedID)
	id := insertResult.InsertedID.(primitive.ObjectID)
	return request, id
}

func FindOwnershipTransferByHex(id string, res http.ResponseWriter) (model.PendingOwnershipTransfer, bool) {
	collection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	requestIdentifier, iderr := primitive.ObjectIDFromHex(id)
	if iderr != nil {
		log.Fatal(iderr)
	}
	filter := bson.M{"_id": requestIdentifier}
	var result model.PendingOwnershipTransfer
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find request!")
			return result, false
		}
		log.Fatal(err)
	}
	return result, true
}
