package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	uri           string
	clientOptions *options.ClientOptions
	client        *mongo.Client
	err           error
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

func SaveUserApplication(email string, firstName string, lastName string, phone string, date string, username string, password string, verified bool, applicationMethod string) {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	// Creating user instance
	user := model.ApplicationUser{
		Email:             email,
		FirstName:         firstName,
		LastName:          lastName,
		Phone:             phone,
		DateOfBirth:       date,
		Username:          username,
		Password:          password,
		Companies:         make([]model.Companies, 0),
		Verified:          verified,
		ApplicationMethod: applicationMethod,
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
	filter := bson.M{"Email": email}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	if CheckPasswordHash(password, result.Password) {
		return true
	}
	return false
}

// ValidEmail checks if the given email exists in the database.
func FindUserEmail(email string) bool {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		fmt.Println(err)
	}
	return true
}

// ValidUsername checks if the given username exists in the database.
func FindUserUsername(username string) bool {
	UsersCollection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Username": username}
	var result model.ApplicationUser
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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

func FindUserByHex(hex string, res http.ResponseWriter) (model.ApplicationUser, bool) {
	collection := GetClient().Database("UserDatabase").Collection("Users")
	userIdentifier, iderr := primitive.ObjectIDFromHex(hex)
	if iderr != nil {
		log.Fatal(iderr)
	}
	filter := bson.M{"_id": userIdentifier}
	var result model.ApplicationUser
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return result, false
		}
	}
	return result, true
}

func FindUserByMail(mail string, res http.ResponseWriter) (model.ApplicationUser, bool) {
	collection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": mail}
	var result model.ApplicationUser
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return result, false
		}
	}
	return result, true
}

func DetermineUsersRoleWithinCompany(user model.ApplicationUser, companyID primitive.ObjectID) string {
	var role string
	for _, company := range user.Companies {
		if company.CompanyID == companyID {
			role = company.Role
		}
	}
	return role
}

func FindUserById(id primitive.ObjectID, res http.ResponseWriter) (model.ApplicationUser, bool) {
	collection := GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"_id": id}
	var result model.ApplicationUser
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return result, false
		}
	}
	return result, true
}
