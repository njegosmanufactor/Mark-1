package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setting up the URL to connect to the MongoDB server
// var Uri, _ = os.LookupEnv("MONGO_URI")
/*var Uri = "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
var ClientOptions = options.Client().ApplyURI(Uri)
var Client, Err = mongo.Connect(context.Background(), ClientOptions)
*/
var (
	uri           string
	clientOptions *options.ClientOptions
	client        *mongo.Client
	err           error
)

// InitConnection initializes the MongoDB connection
func InitConnection() {
	uri, _ = os.LookupEnv("MONGO_URI")
	clientOptions = options.Client().ApplyURI(uri)
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
}

// GetClient returns the MongoDB client for reuse in other packages
func GetClient() *mongo.Client {
	return client
}

// save user into database
func SaveUserOther(email string) {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	// Creating user instance
	user := model.OtherUser{
		Email: email,
	}

	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func SaveUserApplication(email string, firstName string, lastName string, phone string, date string, username string, password string, verified bool) {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	// Creating user instance
	user := model.ApplicationUser{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		Phone:       phone,
		DateOfBirth: date,
		Username:    username,
		Password:    password,
		Company:     "",
		Role:        "User",
		Verified:    verified,
	}

	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func ValidUser(email string, password string) bool {

	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email, "Password": password}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func ValidEmail(email string) bool {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	fmt.Println(UsersCollection.Name() + "mailcontroler")
	filter := bson.M{"Email": email}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Nema ga")
			return false
		}
		fmt.Println(err)
	}
	return true
}

func ValidUsername(username string) bool {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	fmt.Println(UsersCollection.Name())
	filter := bson.M{"Username": username}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	fmt.Println(result.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Nema ni njega")
			return false
		}
		fmt.Println(err)
	}
	return true
}

func VerifyUser(email string) bool {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	update := bson.M{"$set": bson.M{"Verified": true}}
	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	return err == nil
}

func SetUserRoleOwner(userID primitive.ObjectID) error {

	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"Role": "Owner"}}
	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
func GetUserData(email string) (model.ApplicationUser, error) {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return model.ApplicationUser{}, err
	}

	return result, nil
}
