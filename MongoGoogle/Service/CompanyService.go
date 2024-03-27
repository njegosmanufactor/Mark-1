package Service

import (
	model "MongoGoogle/Model"
	dataBase "MongoGoogle/Repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateComapny(res http.ResponseWriter, req *http.Request) {
	user, tokenpointer := GetUserAndPointerFromToken(res, req)
	if tokenpointer != nil && tokenpointer.Valid {
		var companyData struct {
			Name                  string         `json:"name"`
			Address               model.Location `json:"location"`
			Website               string         `json:"website"`
			ListOfApprovedDomains []string       `json:"listOfApprovedDomains"`
		}

		err := json.NewDecoder(req.Body).Decode(&companyData)
		if err != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}

		proba, findCompanyerr := dataBase.FindCompanyByName(companyData.Name, res)
		fmt.Println(proba)
		if findCompanyerr {
			json.NewEncoder(res).Encode("company exist")
		} else {
			dataBase.SaveCompany(companyData.Name, companyData.Address, companyData.Website, companyData.ListOfApprovedDomains, user.ID)
			user, _ := dataBase.GetUserData(user.Email)
			token, _ := GenerateToken(user, time.Hour)
			json.NewEncoder(res).Encode(token)
		}
	} else {
		json.NewEncoder(res).Encode("User not found")
	}
}

func DeleteCompany(res http.ResponseWriter, req *http.Request) {
	user, tokenpointer := GetUserAndPointerFromToken(res, req)

	var requestBody struct {
		CompanyName string `json:"companyName"`
	}
	errReq := json.NewDecoder(req.Body).Decode(&requestBody)
	if errReq != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}
	if requestBody.CompanyName == "" {
		http.Error(res, "Company name is required", http.StatusBadRequest)
		return
	} else {
		company, err := dataBase.FindCompanyByName(requestBody.CompanyName, res)
		if !err {
			json.NewEncoder(res).Encode("You don't have any company")
			return
		}
		if tokenpointer != nil && tokenpointer.Valid && user.ID == company.Owner {
			dataBase.DeleteCompany(requestBody.CompanyName, company.ID)
			user, _ := dataBase.GetUserData(user.Email)
			token, _ := GenerateToken(user, time.Hour)
			json.NewEncoder(res).Encode(token)
		} else {
			json.NewEncoder(res).Encode("You are not owner of " + requestBody.CompanyName)
		}
	}
}

// Function that finds the right company and inserts users id to employees field
func AddUserToCompany(companyId primitive.ObjectID, userEmail string, res http.ResponseWriter) (model.Company, bool) {
	//finding the company
	collection := dataBase.GetClient().Database("UserDatabase").Collection("Company")
	UserCollection := dataBase.GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"_id": companyId}
	var result model.Company
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find company!")
			return result, false
		}
	}
	User := model.Employee{
		Email: userEmail,
		Role:  "User",
	}
	//Updating employees field with the right user
	update := bson.M{"$push": bson.M{"Employees": User}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	userFilter := bson.M{"Email": userEmail}
	Comp := model.Companies{
		CompanyID: companyId,
		Role:      "User",
	}
	update = bson.M{"$push": bson.M{"Companies": Comp}}
	_, err = UserCollection.UpdateOne(context.Background(), userFilter, update)

	if err != nil {
		json.NewEncoder(res).Encode("Company not updated!")
		return result, false
	}
	json.NewEncoder(res).Encode("Company updated!")
	return result, true
}
