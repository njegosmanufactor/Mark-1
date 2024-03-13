package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// save user into database
func SaveCompany(name string, location model.Location, website string, listOfApprovedDomains []string) {
	client, err := Client, Err
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

func DeleteCompany(companyName string) {
	client, err := Client, Err
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Pristupanje kolekciji Company
	companyCollection := client.Database("UserDatabase").Collection("Company")

	// Brisanje kompanije na osnovu njenog imena
	deleteResult, err := companyCollection.DeleteOne(context.Background(), bson.M{"Name": companyName})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted company with name '%s'. Deleted count: %d\n", companyName, deleteResult.DeletedCount)
}

func ValidComapnyName(companyName string) bool {
	client, err := Client, Err
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
