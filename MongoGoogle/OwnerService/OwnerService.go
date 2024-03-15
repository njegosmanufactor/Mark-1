package OwnerService

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	mail "MongoGoogle/ApplicationService"
	conn "MongoGoogle/Repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OwnershipDTO is a data transfer object used for transferring ownership information.
type OwnershipDTO struct {
	OwnerID   string `bson:"_id,omitempty"` //Owner's or admin's ID
	Email     string `bson:"Email"`         //User he sends the invitation to.
	CompanyID string `bson:"CompanyID"`
}

// InvitationDTO is a data transfer object used for transferring invitation information.
type InvitationDTO struct {
	Email     string `bson:"email,omitempty"`
	CompanyID string `bson:"companyId,omitempty"`
}

// Using users id to check his role, we can send invitation to other users for them to join our organisation. Or accept ownership
func TransferOwnership(res http.ResponseWriter, req *http.Request) {

	//Getting request body
	var ownership OwnershipDTO
	decErr := json.NewDecoder(req.Body).Decode(&ownership)
	if decErr != nil {
		http.Error(res, decErr.Error(), http.StatusBadRequest)
	}
	//parsing request data to corect type
	ownerId, iderr := primitive.ObjectIDFromHex(ownership.OwnerID)
	if iderr != nil {
		log.Fatal(iderr)
	}
	companyId, iderr := primitive.ObjectIDFromHex(ownership.CompanyID)
	if iderr != nil {
		log.Fatal(iderr)
	}
	user, found := conn.FindUserByHex(ownership.OwnerID, res)
	if found {
		if user.Role == "Owner" || user.Role == "Admin" {
			//kreiramo zahtev(owner id comp id user mail) ovde i saljemo mejl(id zahteva)
			_, transferId := conn.CreatePendingOwnershipInvitation(ownership.Email, ownerId, companyId)
			mail.SendOwnershipMail(transferId.Hex(), ownership.Email, res)
		} else {
			json.NewEncoder(res).Encode("Only owners can transfer ownership.")
		}
	} else {
		json.NewEncoder(res).Encode("Owner not found!")
	}
}

// FinaliseOwnershipTransfer finalizes the ownership transfer by updating the user's role to "Owner".
func FinaliseOwnershipTransfer(id string, res http.ResponseWriter) {
	userCollection := conn.GetClient().Database("UserDatabase").Collection("Users")
	companyCollection := conn.GetClient().Database("UserDatabase").Collection("Company")
	transferCollection := conn.GetClient().Database("UserDatabase").Collection("PendingRequests")
	transfer, found := conn.FindOwnershipTransferByHex(id, res)
	if found {
		//update usera
		userFilter := bson.M{"Email": transfer.Email}
		updateUser := bson.M{"$set": bson.M{"Role": "Owner"}}
		result, err := userCollection.UpdateOne(context.Background(), userFilter, updateUser)
		if result.MatchedCount == 0 {

			json.NewEncoder(res).Encode("User role not updated")
			log.Fatal(err)
		}
		//update vlasnika
		ownerFilter := bson.M{"_id": transfer.OwnerID}
		updateOwner := bson.M{"$set": bson.M{"Role": "User"}} // When asigning chage role to admin or user?s
		result, err = userCollection.UpdateOne(context.Background(), ownerFilter, updateOwner)
		if result.MatchedCount == 0 {

			json.NewEncoder(res).Encode("Owner role not updated")
			log.Fatal(err)
		}
		//update kompanije
		companyFilter := bson.M{"_id": transfer.CompanyID}
		user, found := conn.FindUserByMail(transfer.Email, res)
		if found {
			updateCompany := bson.M{"$set": bson.M{"Owner": user.ID}} // When asigning chage role to admin or user?s
			result, err = companyCollection.UpdateOne(context.Background(), companyFilter, updateCompany)
			if result.MatchedCount == 0 {

				json.NewEncoder(res).Encode("Owner role not updated")
				log.Fatal(err)
			}
		}
		//update zahteva
		requestFIlter := bson.M{"_id": transfer.ID}
		updateRequest := bson.M{"$set": bson.M{"Completed": true}} // When asigning chage role to admin or user?s
		result, err = transferCollection.UpdateOne(context.Background(), requestFIlter, updateRequest)
		if result.MatchedCount == 0 {

			json.NewEncoder(res).Encode("Owner role not updated")
			log.Fatal(err)
		}

	} else {
		json.NewEncoder(res).Encode("Transfer not found")

	}
}

// Sends an invitation to join a company to the specified email address.
func SendInvitation(res http.ResponseWriter, req *http.Request) {
	var invitation InvitationDTO
	decErr := json.NewDecoder(req.Body).Decode(&invitation)
	if decErr != nil {
		http.Error(res, decErr.Error(), http.StatusBadRequest)
	}
	//finding the user
	user, found := conn.FindUserByMail(invitation.Email, res)
	if found {
		if user.Verified {
			_, id := conn.CreatePendingInvite(invitation.Email, invitation.CompanyID)
			mail.SendInvitationMail(id.Hex(), invitation.Email)
		} else {
			json.NewEncoder(res).Encode("This user hasn't verified his account.")
		}
	}
}
