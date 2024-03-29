package Service

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
		appUser, err := data.GetUserData(user.Email)
		if err != nil {
			json.NewEncoder(res).Encode(err)
		} else {
			if appUser.ApplicationMethod != "Google" {
				json.NewEncoder(res).Encode(user.Email + " already exists")
			}
		}
	} else {
		data.SaveUserApplication(user.Email, user.GivenName, user.FamilyName, "", "", user.Email, "", true, "Google")
		//check if user has invites prior to registering
		//pending repo that returns the invite id
		id, found := CheckForUnregInvites(user.Email, res)
		//mail service that sends the invite
		if found {
			SendInvitationMail(id.Hex(), user.Email)
		}

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
