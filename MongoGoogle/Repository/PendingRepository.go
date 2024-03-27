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

// Creates a password change request in the database and returns the created request and its ID.
func CreatePasswordChangeRequest(email string) (model.PasswordChangeRequest, primitive.ObjectID, bool) {
	RequestCollection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	request := model.PasswordChangeRequest{
		Email:     email,
		Completed: false,
	}
	var res http.ResponseWriter
	insertResult, err := RequestCollection.InsertOne(context.Background(), request)
	if err != nil {
		json.NewEncoder(res).Encode(err)
	} else {
		fmt.Println("Added new request with ID:", insertResult.InsertedID)
		id := insertResult.InsertedID.(primitive.ObjectID)
		return request, id, true
	}
	var id primitive.ObjectID
	return request, id, false
}

// Creates a pending invite request in the database and returns the created request and its ID.
func CreatePendingInvite(email string, companyId string) (model.PendingRequest, primitive.ObjectID, bool) {
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
	var res http.ResponseWriter
	// Adding request to the database
	insertResult, err := RequestCollection.InsertOne(context.Background(), request)
	if err != nil {
		json.NewEncoder(res).Encode(err)
	} else {
		fmt.Println("Added new request with ID:", insertResult.InsertedID)
		id := insertResult.InsertedID.(primitive.ObjectID)
		return request, id, true
	}
	var id primitive.ObjectID
	return request, id, false
}

// Creates a passwordless request in the database and returns the created request and its code.
func CreatePasswordLessRequest(email string, code string) bool {
	RequestCollection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	// Creating request instance
	request := model.PasswordLessRequest{
		Email:     email,
		Code:      code,
		Completed: false,
	}
	var res http.ResponseWriter
	// Adding request to the database
	RequestCollection.InsertOne(context.Background(), request)
	if err != nil {
		json.NewEncoder(res).Encode(err)
	} else {
		fmt.Println("CODE: ", code)
		return true
	}
	return false
}

// Creates a pending ownership invitation request in the database and returns the created request and its ID.
func CreatePendingOwnershipInvitation(email string, ownerId primitive.ObjectID, companyId primitive.ObjectID) (model.PendingOwnershipTransfer, primitive.ObjectID, bool) {
	RequestCollection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	// Creating request instance
	request := model.PendingOwnershipTransfer{
		Email:     email,
		CompanyID: companyId,
		OwnerID:   ownerId,
		Completed: false,
	}
	// Adding request to the database
	var res http.ResponseWriter
	insertResult, err := RequestCollection.InsertOne(context.Background(), request)
	if err != nil {
		json.NewEncoder(res).Encode(err)
	} else {
		fmt.Println("Added new request with ID:", insertResult.InsertedID)
		id := insertResult.InsertedID.(primitive.ObjectID)
		return request, id, true
	}
	var id primitive.ObjectID
	return request, id, false
}

// Finds an ownership transfer request by its ID in the database and returns the request and a boolean indicating if it was found.
func FindOwnershipTransferByHex(id string, res http.ResponseWriter) (model.PendingOwnershipTransfer, bool) {
	collection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	requestIdentifier, iderr := primitive.ObjectIDFromHex(id)
	if iderr != nil {
		json.NewEncoder(res).Encode(iderr)
	}
	filter := bson.M{"_id": requestIdentifier}
	var result model.PendingOwnershipTransfer
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find request!")
			return result, false
		}
	}
	return result, true
}

// Finds a passwordless request by its ID in the database and returns the request and a boolean indicating if it was found.
func FindCodeRequestByHex(id string, res http.ResponseWriter) (model.PasswordLessRequest, bool) {
	collection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	requestIdentifier, iderr := primitive.ObjectIDFromHex(id)
	if iderr != nil {
		json.NewEncoder(res).Encode(iderr)
	}
	filter := bson.M{"_id": requestIdentifier}
	var result model.PasswordLessRequest
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find request!")
			return result, false
		}
	}
	return result, true
}

// Deletes a pending request from the database based on its ID.
func DeletePandingRequrst(id string) {
	collection := GetClient().Database("UserDatabase").Collection("PendingRequests")
	requestIdentifier, iderr := primitive.ObjectIDFromHex(id)
	var res http.ResponseWriter
	if iderr != nil {
		json.NewEncoder(res).Encode(iderr)
	} else {
		filter := bson.M{"_id": requestIdentifier}
		collection.DeleteOne(context.Background(), filter)
	}
}
