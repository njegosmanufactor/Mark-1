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
