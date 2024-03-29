package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Saves a new company into the database.
func SaveCompany(name string, location model.Location, website string, listOfApprovedDomains []string, ownerId primitive.ObjectID) {
	CompanyCollection := GetClient().Database("UserDatabase").Collection("Company")
	UserCollection := GetClient().Database("UserDatabase").Collection("Users")

	// Creating user instance
	company := model.Company{
		Name:                  name,
		Address:               location,
		Website:               website,
		ListOfApprovedDomains: listOfApprovedDomains,
		Owner:                 ownerId,
		Employees:             make([]model.Employee, 0),
	}
	var res http.ResponseWriter
	user, _ := FindUserById(ownerId, res)
	Employee := model.Employee{
		Email: user.Email,
		Role:  "Owner",
	}
	company.Employees = append(company.Employees, Employee)
	// Adding user to the database
	insertResult, err := CompanyCollection.InsertOne(context.Background(), company)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}
	filter := bson.M{"Email": user.Email}
	Comp := model.Companies{
		CompanyID: insertResult.InsertedID.(primitive.ObjectID),
		Role:      "Owner",
	}
	update := bson.M{"$push": bson.M{"Companies": Comp}}

	_, err = UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}
	fmt.Println("Added new company with ID:", insertResult.InsertedID)
}

// Deletes a company from the database based on its name.
func DeleteCompany(companyName string, companyID primitive.ObjectID) {
	companyCollection := GetClient().Database("UserDatabase").Collection("Company")
	var res http.ResponseWriter
	company, Comperr := FindCompanyByName(companyName, res)
	if Comperr {
		for _, Employee := range company.Employees {
			RemoveCompanyFromUser(Employee.Email, company.ID)
		}
		deleteResult, err := companyCollection.DeleteOne(context.Background(), bson.M{"Name": companyName})
		if err != nil {
			json.NewEncoder(res).Encode(err)
			return
		}
		fmt.Printf("Deleted company with name ' %s'. Deleted count: %d\n", companyName, deleteResult.DeletedCount)
	} else {
		json.NewEncoder(res).Encode("Company not found!")
	}

}
func RemoveCompanyFromUser(mail string, companyId primitive.ObjectID) {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")

	filter := bson.M{"Email": mail}
	update := bson.M{"$pull": bson.M{"Companies": bson.M{"_id": companyId}}}

	_, err := UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		var res http.ResponseWriter
		json.NewEncoder(res).Encode(err)
	}

}

// Returns pair (company,bool). True if the company is found, false if not
func FindCompanyByHex(companyId string, res http.ResponseWriter) (model.Company, bool) {
	collection := GetClient().Database("UserDatabase").Collection("Company")
	userIdentifier, iderr := primitive.ObjectIDFromHex(companyId)
	if iderr != nil {
		var res http.ResponseWriter
		json.NewEncoder(res).Encode(iderr)

	}
	filter := bson.M{"_id": userIdentifier}
	var result model.Company
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find company!")
			return result, false
		}
	}
	return result, true
}

// Finds a company by its name in the database.
func FindCompanyByName(companyName string, res http.ResponseWriter) (model.Company, bool) {
	collection := GetClient().Database("UserDatabase").Collection("Company")
	filter := bson.M{"Name": companyName}
	var result model.Company
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, false
		}
	}
	return result, true
}
