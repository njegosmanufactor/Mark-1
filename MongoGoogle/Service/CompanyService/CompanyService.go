package CompanyService

import (
	model "MongoGoogle/Model"
	dataBase "MongoGoogle/Repository"
	tokenService "MongoGoogle/Service/TokenService"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CreateComapny(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	token = tokenService.SplitTokenHeder(token)
	user, tokenUser, err := tokenService.ExtractUserFromToken(token)
	if err != nil {
		http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
		return
	}
	if tokenUser != nil && tokenUser.Valid {
		var companyData struct {
			Name                  string         `json:"name"`
			Address               model.Location `json:"location"`
			Website               string         `json:"website"`
			ListOfApprovedDomains []string       `json:"listOfApprovedDomains"`
		}

		err := json.NewDecoder(req.Body).Decode(&companyData)
		if err != nil {
			http.Error(res, "Error decoding request body", http.StatusBadRequest)
			return
		}

		if dataBase.FindComapnyName(companyData.Name) {
			fmt.Printf("Company exist\n")
		} else {
			dataBase.SaveCompany(companyData.Name, companyData.Address, companyData.Website, companyData.ListOfApprovedDomains, user.ID)
			user, _ := dataBase.GetUserData(user.Email)
			token, _ := tokenService.GenerateToken(user, time.Hour)
			json.NewEncoder(res).Encode(token)
		}
	} else {
		json.NewEncoder(res).Encode("User not found")
	}
}

func DeleteCompany(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	token = tokenService.SplitTokenHeder(token)
	user, tokenUser, err := tokenService.ExtractUserFromToken(token)
	if err != nil {
		http.Error(res, "Error extracting user from token", http.StatusInternalServerError)
		return
	}

	var requestBody struct {
		CompanyName string `json:"companyName"`
	}
	errReq := json.NewDecoder(req.Body).Decode(&requestBody)
	if errReq != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}
	if requestBody.CompanyName == "" {
		http.Error(res, "Company name is required", http.StatusBadRequest)
		return
	} else {
		company, err := dataBase.FindCompanyByName(requestBody.CompanyName, res)
		if !err {
			json.NewEncoder(res).Encode("You don't have any company")
			return
		}
		if tokenUser != nil && tokenUser.Valid && user.ID == company.Owner {
			dataBase.DeleteCompany(requestBody.CompanyName, company.ID)
			user, _ := dataBase.GetUserData(user.Email)
			token, _ := tokenService.GenerateToken(user, time.Hour)
			json.NewEncoder(res).Encode(token)
		} else {
			json.NewEncoder(res).Encode("You are not owner of" + requestBody.CompanyName)
		}
	}
}
