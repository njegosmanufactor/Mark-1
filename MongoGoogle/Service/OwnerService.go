package Service

import (
	model "MongoGoogle/Model"
	conn "MongoGoogle/Repository"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		json.NewEncoder(res).Encode(iderr)
		return
	}
	companyId, iderr := primitive.ObjectIDFromHex(ownership.CompanyID)
	if iderr != nil {
		json.NewEncoder(res).Encode(iderr)
		return
	}
	user, found := conn.FindUserByHex(ownership.OwnerID, res)
	role := conn.DetermineUsersRoleWithinCompany(user, companyId)
	if found {
		if role == "Owner" {
			//kreiramo zahtev(owner id comp id user mail) ovde i saljemo mejl(id zahteva)
			_, transferId, created := conn.CreatePendingOwnershipInvitation(ownership.Email, ownerId, companyId)
			if created {
				SendOwnershipMail(transferId.Hex(), ownership.Email, res)
			} else {
				json.NewEncoder(res).Encode("Error on creating ownership invitation")
			}
		} else {
			json.NewEncoder(res).Encode("Only owners can transfer ownership.")
		}
	} else {
		json.NewEncoder(res).Encode("Owner not found!")
	}
}
func UpdateUserRole(transfer model.PendingOwnershipTransfer, res http.ResponseWriter, userCollection *mongo.Collection) bool {
	//update usera                                    //prolazi ako je korisnik vec user u firmi
	userFilter := bson.M{"Email": transfer.Email, "Companies._id": transfer.CompanyID}
	updateUser := bson.M{"$set": bson.M{"Companies.$.Role": "Owner"}}
	result, err := userCollection.UpdateOne(context.Background(), userFilter, updateUser)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return false
	}
	if result.MatchedCount == 0 {
		json.NewEncoder(res).Encode("User role not updated")
		return false
	}
	return true
}
func UpdateOwnerRole(transfer model.PendingOwnershipTransfer, res http.ResponseWriter, userCollection *mongo.Collection) bool {
	ownerFilter := bson.M{"_id": transfer.OwnerID, "Companies._id": transfer.CompanyID}
	updateOwner := bson.M{"$set": bson.M{"Companies.$.Role": "User"}} // When asigning chage role to admin or user?s
	result, err := userCollection.UpdateOne(context.Background(), ownerFilter, updateOwner)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return false
	}
	if result.MatchedCount == 0 {
		json.NewEncoder(res).Encode("User role not updated")
		return false
	}
	return true
}

func UpdateCompanyEmployees(transfer model.PendingOwnershipTransfer, res http.ResponseWriter, companyCollection *mongo.Collection) bool {
	companyFilter := bson.M{"_id": transfer.CompanyID}
	user, found := conn.FindUserByMail(transfer.Email, res)
	if found {
		updateCompany := bson.M{"$set": bson.M{"OwnerId": user.ID}} // When asigning chage role to admin or user?s
		result, err := companyCollection.UpdateOne(context.Background(), companyFilter, updateCompany)
		if err != nil {
			json.NewEncoder(res).Encode(err)
			return false
		}
		if result.MatchedCount == 0 {

			json.NewEncoder(res).Encode("Owner role not updated")
			return false
		}
		companyEmployeeFilter := bson.M{"_id": transfer.CompanyID, "Employees.Email": transfer.Email}
		updateCompanyOwner := bson.M{"$set": bson.M{"Employees.$.Role": "Owner"}}
		result, err = companyCollection.UpdateOne(context.Background(), companyEmployeeFilter, updateCompanyOwner)
		if err != nil {
			json.NewEncoder(res).Encode(err)
			return false
		}
		if result.MatchedCount == 0 {

			json.NewEncoder(res).Encode("Owner role not updated")
			return false
		}
		owner, found := conn.FindUserById(transfer.OwnerID, res)
		if found {
			companyEmployeeOwnerFilter := bson.M{"_id": transfer.CompanyID, "Employees.Email": owner.Email}
			updateCompanyEmployee := bson.M{"$set": bson.M{"Employees.$.Role": "User"}}
			result, err = companyCollection.UpdateOne(context.Background(), companyEmployeeOwnerFilter, updateCompanyEmployee)
			if err != nil {
				json.NewEncoder(res).Encode(err)
				return false
			}
			if result.MatchedCount == 0 {

				json.NewEncoder(res).Encode("Owner role not updated")
				return false
			}
		}
	}
	return true
}

// FinaliseOwnershipTransfer finalizes the ownership transfer by updating the user's role to "Owner".
func FinaliseOwnershipTransfer(id string, res http.ResponseWriter) {
	userCollection := conn.GetClient().Database("UserDatabase").Collection("Users")
	companyCollection := conn.GetClient().Database("UserDatabase").Collection("Company")
	transferCollection := conn.GetClient().Database("UserDatabase").Collection("PendingRequests")
	transfer, found := conn.FindOwnershipTransferByHex(id, res)
	if found {
		//update usera                                    //prolazi ako je korisnik vec user u firmi
		if !UpdateUserRole(transfer, res, userCollection) {
			json.NewEncoder(res).Encode("User role not updated")
			return
		}

		//update vlasnika
		if !UpdateOwnerRole(transfer, res, userCollection) {
			json.NewEncoder(res).Encode("User role not updated")
			return
		}
		//update kompanije
		if !UpdateCompanyEmployees(transfer, res, companyCollection) {
			json.NewEncoder(res).Encode("Company not updated")
			return
		}
		//update zahteva
		requestFIlter := bson.M{"_id": transfer.ID}
		updateRequest := bson.M{"$set": bson.M{"Completed": true}} // When asigning chage role to admin or user?s
		result, err := transferCollection.UpdateOne(context.Background(), requestFIlter, updateRequest)
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
			_, id, created := conn.CreatePendingInvite(invitation.Email, invitation.CompanyID)
			if created {
				SendInvitationMail(id.Hex(), invitation.Email)
			} else {
				json.NewEncoder(res).Encode("Error on creating pending invite")
			}
		} else {
			json.NewEncoder(res).Encode("This user hasn't verified his account.")
		}
	} else {
		// Didnt find the user specifiend in body. I should create pending invite, and when user registers there should be a check
		// System goes through requests and if it finds an invite , it sends that invite.
		_, id, created := conn.CreateUnregInvite(invitation.Email, invitation.CompanyID)
		if created {
			json.NewEncoder(res).Encode("Created unreg invite with id:" + id.Hex())
		}
		json.NewEncoder(res).Encode("User must register with provided mail in order to join the company.")
	}
}
