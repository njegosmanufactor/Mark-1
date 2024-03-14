package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Saves a new company into the database.
func SaveCompany(name string, location model.Location, website string, listOfApprovedDomains []string) {
	client = GetClient()
	CompanyCollection := client.Database("UserDatabase").Collection("Company")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	// Creating user instance
	company := model.Company{
		Name:                  name,
		Address:               location,
		Website:               website,
		ListOfApprovedDomains: listOfApprovedDomains,
	}

	// Adding user to the database
	insertResult, err := CompanyCollection.InsertOne(context.Background(), company)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new company with ID:", insertResult.InsertedID)
}

// Deletes a company from the database based on its name.
func DeleteCompany(companyName string) {
	client = GetClient()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	companyCollection := client.Database("UserDatabase").Collection("Company")

	deleteResult, err := companyCollection.DeleteOne(context.Background(), bson.M{"Name": companyName})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted company with name '%s'. Deleted count: %d\n", companyName, deleteResult.DeletedCount)
}

// Checks if a company with the given name exists in the database.
func ValidComapnyName(companyName string) bool {
	client = GetClient()
	UsersCollection := client.Database("UserDatabase").Collection("Company")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"Name": companyName}
	var result model.Company
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

// Sets the company for a user identified by userID in the database.
func SetUserCompany(userID primitive.ObjectID, companyName string) error {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"Company": companyName}}

	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
