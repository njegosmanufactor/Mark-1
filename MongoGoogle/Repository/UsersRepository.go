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
		Authorised:  false,
	}

	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

// ValidUser checks if the user with the given email and password exists in the database.
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

// ValidEmail checks if the given email exists in the database.
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

// ValidUsername checks if the given username exists in the database.
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

// VerifyUser updates the verification status of the user with the given email.
func VerifyUser(email string) bool {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	update := bson.M{"$set": bson.M{"Verified": true}}
	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	return err == nil
}

// SetUserRole updates the role of the user with the given user ID.
func SetUserRole(userID primitive.ObjectID, role string) error {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"Role": role}}

	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// GetUserData retrieves the user data for the given email.
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

func SetAuthorise(userID primitive.ObjectID, authorised bool) error {
	client, err := MongoConnection()
	UsersCollection := client.Database("UserDatabase").Collection("Users")
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"Authorised": authorised}}

	_, err = UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
