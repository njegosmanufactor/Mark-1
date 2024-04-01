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

func InsertInflow(Name string, CashflowID string, UserID string, Duration string, Amount float32, Category string, Subcategory string, res http.ResponseWriter) bool {
	CashFlowCollection := GetClient().Database("UserDatabase").Collection("CashFlow")
	cashflowID, cfErr := primitive.ObjectIDFromHex(CashflowID)
	if cfErr != nil {
		json.NewEncoder(res).Encode(cfErr)
	}
	userID, usErr := primitive.ObjectIDFromHex(UserID)
	if cfErr != nil {
		json.NewEncoder(res).Encode(usErr)
	}
	transaciton := model.Transaction{
		Name:     Name,
		UserID:   userID,
		Duration: Duration,
		Amount:   Amount,
	}
	filter := bson.M{"_id": cashflowID, "Categories": bson.M{"$elemMatch": bson.M{"Name": Category, "Subcategories": bson.M{"$elemMatch": bson.M{"Name": Subcategory}}}}}
	var result model.CashFlow
	err := CashFlowCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			categoryFilter := bson.M{"_id": cashflowID, "Categories": bson.M{"$elemMatch": bson.M{"Name": Category}}}
			subcategory := model.Subcategory{
				Name:    Subcategory,
				Inflow:  make([]model.Transaction, 0),
				Outflow: make([]model.Transaction, 0),
			}
			subcategory.Inflow = append(subcategory.Inflow, transaciton)
			update := bson.M{"$push": bson.M{"Categories.$.Subcategories": subcategory}}
			_, err = CashFlowCollection.UpdateOne(context.Background(), categoryFilter, update)
			if err != nil {
				fmt.Println(err)
				return false
			}
		}
	}
	subcategoryFilter := bson.M{"_id": cashflowID, "Categories": bson.M{"$elemMatch": bson.M{"Name": Category}}}

	update := bson.M{"$push": bson.M{"Categories.2.Subcategories.$.Inflow": transaciton}}
	_, err = CashFlowCollection.UpdateOne(context.Background(), subcategoryFilter, update)
	if err != nil {
		fmt.Println(err)
		return false
	}
	json.NewEncoder(res).Encode(result)
	return true
}
