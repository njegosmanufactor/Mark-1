package OwnerService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mail "MongoGoogle/ApplicationService"
	conn "MongoGoogle/MongoDB"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OwnershipDTO struct {
	ID    string `bson:"_id,omitempty"` //Owner's or admin's ID
	Email string `bson:"Email"`         //User he sends the invitation to.
}
type InvitationDTO struct {
	Email string `bson:"Email"`
	ID    string `bson:"_id,omitempty"`
}

// Using users id to check his role, we can send invitation to other users for them to join our organisation. Or accept ownership
func TransferOwnership(res http.ResponseWriter, req *http.Request) {

	//Getting request body
	var ownership OwnershipDTO
	decErr := json.NewDecoder(req.Body).Decode(&ownership)
	if decErr != nil {
		http.Error(res, decErr.Error(), http.StatusBadRequest)
	}

	//Finding user in database
	collection := conn.Client.Database("UserDatabase").Collection("Users")
	identifier, iderr := primitive.ObjectIDFromHex(ownership.ID)
	if iderr != nil {
		log.Fatal(iderr)
	}
	filter := bson.M{"_id": identifier}
	fmt.Println(identifier)
	var result conn.ApplicationUser
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
		}
	}

	if result.Role == "Owner" || result.Role == "Admin" {
		//It is assumed the user is already in company. On front, the admin will only have list of those kind of users.
		mail.SendOwnershipMail(ownership.Email, res)
	} else {
		json.NewEncoder(res).Encode("Only owners can transfer ownership.")
	}

}
func FinaliseOwnershipTransfer(email string) error {
	collection := conn.Client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	update := bson.M{"$set": bson.M{"Role": "Owner"}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func SendInvitation(res http.ResponseWriter, req *http.Request) {

	var invitation InvitationDTO
	decErr := json.NewDecoder(req.Body).Decode(&invitation)
	if decErr != nil {
		http.Error(res, decErr.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(res).Encode(invitation)

	//finding the user
	collection := conn.Client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": invitation.Email}
	var result conn.ApplicationUser
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return
		}
	} else { //this is company id extracted from admins or owners profile
		mail.SendInvitationMail(invitation.Email, invitation.ID)
	}

}
