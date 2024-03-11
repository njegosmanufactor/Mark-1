package MongoDB

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OtherUser struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Email string             `bson:"Email"`
}

type ApplicationUser struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Email       string             `bson:"Email"`
	FirstName   string             `bson:"FirstName"`
	LastName    string             `bson:"LastName"`
	Phone       string             `bson:"Phone"`
	DateOfBirth string             `bson:"DateOfBirth"`
	Username    string             `bson:"Username"`
	Password    string             `bson:"Password"`
	Company     string             `bson:"Company"`
	Country     string             `bson:"Country"`
	City        string             `bson:"City"`
	Address     string             `bson:"Address"`
	Role        string             `bson:"Role"`
	Verified    bool               `bson:"Verified"`
}

// save user into database
func SaveUserOther(email string) {
	// Setting up the URL to connect to the MongoDB server
	uri := "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"

	// Setting up client options for connection
	clientOptions := options.Client().ApplyURI(uri)

	// Connecting to the MongoDB server
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Connection to MongoDB successful")
	collection := client.Database("UserDatabase").Collection("Users")

	// Creating user instance
	user := OtherUser{
		Email: email,
	}

	// Adding user to the database
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func SaveUserApplication(email string, firstName string, lastName string, phone string, date string, username string, password string, company string, country string, city string, address string) {
	// Setting up the URL to connect to the MongoDB server
	uri := "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"

	// Setting up client options for connection
	clientOptions := options.Client().ApplyURI(uri)

	// Connecting to the MongoDB server
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Connection to MongoDB successful")
	collection := client.Database("UserDatabase").Collection("Users")

	// Creating user instance
	user := ApplicationUser{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		Phone:       phone,
		DateOfBirth: date,
		Username:    username,
		Password:    password,
		Company:     company,
		Country:     country,
		City:        city,
		Address:     address,
		Role:        "User",
	}

	// Adding user to the database
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
}

func ValidUser(email string, password string) bool {
	uri := "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email, "Password": password}
	var result ApplicationUser
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func ValidEmail(email string) bool {
	uri := "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	var result ApplicationUser
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func ValidUsername(username string) bool {
	uri := "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Username": username}
	var result ApplicationUser
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func SetUserRoleOwner(email string) error {
	uri := "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Povezivanje sa kolekcijom "Users"
	collection := client.Database("UserDatabase").Collection("Users")

	// Kreiranje filtera koji odgovara korisniku sa datim korisničkim imenom
	filter := bson.M{"Email": email}

	// Kreiranje novih vrednosti koje želimo da ažuriramo
	update := bson.M{"$set": bson.M{"Role": "Owner"}}

	// Ažuriranje dokumenta u bazi
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func GetUserData(email string) (ApplicationUser, error) {
	// Povezivanje sa MongoDB bazom
	uri := "mongodb+srv://Nikola045:Bombarder535@userdatabase.qcrmscd.mongodb.net/?retryWrites=true&w=majority&appName=UserDataBase"
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return ApplicationUser{}, err
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Povezivanje sa kolekcijom "Users"
	collection := client.Database("UserDatabase").Collection("Users")

	// Kreiranje filtera koji odgovara korisniku sa datim mejlom
	filter := bson.M{"Email": email}

	// Definisanje strukture za rezultat
	var result ApplicationUser

	// Dohvatanje korisničkih podataka iz baze
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return ApplicationUser{}, err
	}

	return result, nil
}
