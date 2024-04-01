package Service

import (
	dataBase "MongoGoogle/Repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CreateCashFlowUser(res http.ResponseWriter, req *http.Request) {
	user, tokenpointer := GetUserAndPointerFromToken(res, req)
	if tokenpointer != nil && tokenpointer.Valid {
		dataBase.CreateChashFlowForUser(user.ID)
		user, _ := dataBase.GetUserData(user.Email)
		token, _ := GenerateToken(user, time.Hour)
		json.NewEncoder(res).Encode(token)
	} else {
		json.NewEncoder(res).Encode("User not found")
	}
}

func CreateCashFlowCompany(res http.ResponseWriter, req *http.Request) {
	user, tokenpointer := GetUserAndPointerFromToken(res, req)
	fmt.Println(user)
	if tokenpointer != nil && tokenpointer.Valid {
		dataBase.CreateChashFlowForCompany(user.Companies[0].CompanyID)
		json.NewEncoder(res).Encode("Cash Flow created")
	} else {
		json.NewEncoder(res).Encode("User not found")
	}
}

func CreateInflow(res http.ResponseWriter, req *http.Request) {
	_, tokenpointer := GetUserAndPointerFromToken(res, req)
	var request struct {
		Name        string  `json:"name"`
		CashflowID  string  `json:"cashflowID"`
		UserID      string  `json:"userID"`
		Duration    string  `json:"duration"`
		Amount      float32 `json:"amount"`
		Category    string  `json:"category"`
		Subcategory string  `json:"subcategory"`
	}
	if tokenpointer != nil && tokenpointer.Valid {
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}
		if dataBase.InsertInflow(request.Name, request.CashflowID, request.UserID, request.Duration, request.Amount, request.Category, request.Subcategory, res) {
			json.NewEncoder(res).Encode("Inflow created")
		} else {
			json.NewEncoder(res).Encode("Error on creating inflow.")
			return
		}
	} else {
		json.NewEncoder(res).Encode("User not found")
	}
}
