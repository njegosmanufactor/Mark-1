package GoogleService

import (
	"fmt"
	"net/http"

	data "MongoGoogle/Repository"

	userType "MongoGoogle/Model"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// Completes the user authentication process using Google OAuth.
func CompleteGoogleUserAuthentication(res http.ResponseWriter, req *http.Request) {
	tmpUser, err := gothic.CompleteUserAuth(res, req)
	user := AddUserRole(&tmpUser)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}

	if data.ValidEmail(user.Email) {
		// if user have account
	} else {
		fmt.Println("Account created google")
		data.SaveUserApplication(user.Email, user.FirstName, user.LastName, "", "", user.Email, "", true)
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
