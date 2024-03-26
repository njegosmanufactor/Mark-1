package GoogleService

import (
	"encoding/json"
	"net/http"

	userType "MongoGoogle/Model"
	data "MongoGoogle/Repository"

	oauth2v2 "google.golang.org/api/oauth2/v2"

	"github.com/markbates/goth"
)

// Completes the user authentication process using Google OAuth.
func CompleteGoogleUserAuthentication(res http.ResponseWriter, req *http.Request, user *oauth2v2.Userinfo) {
	if data.FindUserEmail(user.Email) {
		appUser, _ := data.GetUserData(user.Email)
		if appUser.ApplicationMethod != "Google" {
			json.NewEncoder(res).Encode(user.Email + " already exists")
		} else {
			json.NewEncoder(res).Encode(user.Email + " successfully logged in to mark-1")
		}
	} else {
		json.NewEncoder(res).Encode(user.Email + " successfully registred to Mark-1")
		data.SaveUserApplication(user.Email, user.GivenName, user.FamilyName, "", "", user.Email, "", true, "Google")
	}
}

// Adds user role to the Google user data.
func AddUserRole(user *goth.User) userType.GoogleData {
	var roleUser userType.GoogleData

	roleUser.AccessToken = user.AccessToken
	roleUser.AccessTokenSecret = user.AccessTokenSecret
	roleUser.AvatarURL = user.AvatarURL
	roleUser.Description = user.Description
	roleUser.Email = user.Email
	roleUser.ExpiresAt = user.ExpiresAt
	roleUser.FirstName = user.FirstName
	roleUser.IDToken = user.IDToken
	roleUser.LastName = user.LastName
	roleUser.Location = user.Location
	roleUser.Name = user.Name
	roleUser.NickName = user.NickName
	roleUser.Provider = user.Provider
	roleUser.RawData = user.RawData
	roleUser.RefreshToken = user.RefreshToken
	roleUser.UserID = user.UserID
	roleUser.Role = "User"

	return roleUser
}
