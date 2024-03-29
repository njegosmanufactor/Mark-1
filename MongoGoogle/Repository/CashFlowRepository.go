package Repository

import (
	model "MongoGoogle/Model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateChashFlowForUser(userId primitive.ObjectID) {
	CashFlowCollection := GetClient().Database("UserDatabase").Collection("CashFlow")
	UserCollection := GetClient().Database("UserDatabase").Collection("Users")

	CategoryOperating := model.Category{
		Name:          "Operating",
		Subcategories: make([]model.Subcategory, 0),
	}
	CategoryInvestment := model.Category{
		Name:          "Investment",
		Subcategories: make([]model.Subcategory, 0),
	}
	CategoryFinancing := model.Category{
		Name:          "Financing",
		Subcategories: make([]model.Subcategory, 0),
	}

	var categories []model.Category
	categories = append(categories, CategoryOperating)
	categories = append(categories, CategoryInvestment)
	categories = append(categories, CategoryFinancing)
	cashFlow := model.CashFlow{
		Categories: categories,
	}
	var res http.ResponseWriter
	user, _ := FindUserById(userId, res)
	insertResult, err := CashFlowCollection.InsertOne(context.Background(), cashFlow)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}
	filter := bson.M{"Email": user.Email}
	update := bson.M{"$set": bson.M{"CashFlowID": insertResult.InsertedID.(primitive.ObjectID)}}

	_, err = UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}
	fmt.Println("Added new cash flow with ID:", insertResult.InsertedID)
}

func CreateChashFlowForCompany(companyID primitive.ObjectID) {
	CashFlowCollection := GetClient().Database("UserDatabase").Collection("CashFlow")
	CompanyCollection := GetClient().Database("UserDatabase").Collection("Company")

	CategoryOperating := model.Category{
		Name:          "Operating",
		Subcategories: make([]model.Subcategory, 0),
	}
	CategoryInvestment := model.Category{
		Name:          "Investment",
		Subcategories: make([]model.Subcategory, 0),
	}
	CategoryFinancing := model.Category{
		Name:          "Financing",
		Subcategories: make([]model.Subcategory, 0),
	}

	var categories []model.Category
	categories = append(categories, CategoryOperating)
	categories = append(categories, CategoryInvestment)
	categories = append(categories, CategoryFinancing)
	cashFlow := model.CashFlow{
		Categories: categories,
	}
	var res http.ResponseWriter
	compnay, _ := FindCompanyByHex(companyID.Hex(), res)
	insertResult, err := CashFlowCollection.InsertOne(context.Background(), cashFlow)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}
	filter := bson.M{"Name": compnay.Name}
	update := bson.M{"$set": bson.M{"CashFlowID": insertResult.InsertedID.(primitive.ObjectID)}}

	_, err = CompanyCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		json.NewEncoder(res).Encode(err)
		return
	}
	fmt.Println("Added new cash flow with ID:", insertResult.InsertedID)
}
