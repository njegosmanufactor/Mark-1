package Service

import (
	model "MongoGoogle/Model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db "MongoGoogle/Repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	oauth2v2 "google.golang.org/api/oauth2/v2"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

var jwtKey = []byte("tajna_lozinka")

// Generates a JWT token for the given user with a specified expiration time.
func GenerateToken(user model.ApplicationUser, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":                user.ID,
		"email":             user.Email,
		"firstName":         user.FirstName,
		"lastName":          user.LastName,
		"phone":             user.Phone,
		"dateOfBirth":       user.DateOfBirth,
		"username":          user.Username,
		"password":          user.Password,
		"verified":          user.Verified,
		"applicationMethod": user.ApplicationMethod,
		"exp":               time.Now().Add(exp).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// SplitTokenHeader extracts the token from the authorization header.
func SplitTokenHeder(authHeader string) string {

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	tokenString := parts[1]
	return tokenString
}

// Parses the JWT token string and returns the token object.
func ParseTokenString1(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return token, fmt.Errorf("failed to parse token: %v", err)
	}
	return token, nil
}

func GetUserAndPointerFromToken(res http.ResponseWriter, req *http.Request) (model.ApplicationUser, *jwt.Token) {
	token := req.Header.Get("Authorization")
	token = SplitTokenHeder(token)
	tokenpointer, _ := ParseTokenString1(token)
	if tokenpointer == nil {
		var user model.ApplicationUser
		return user, nil
	}
	claims, _ := tokenpointer.Claims.(jwt.MapClaims)
	userID, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	user, _ := db.FindUserById(userID, res)
	return user, tokenpointer
}

// Verifies if the token pointer exists and is not nil.
func VerifyTokenPointer(token *jwt.Token) bool {
	if token == nil {
		return false
	} else {
		return true
	}
}

// Logic for handling user login via the application, including token verification and generation of a new one if needed.
func TokenAppLoginLogic(res http.ResponseWriter, req *http.Request, authHeader string, email string, password string) {
	_, tokenPointer := GetUserAndPointerFromToken(res, req)
	message := ApplicationLogin(email, password)

	if tokenPointer == nil {
		if message == "Success" {
			user, _ := db.GetUserData(email)
			token, _ := GenerateToken(user, time.Hour)
			json.NewEncoder(res).Encode(token)
		} else {
			json.NewEncoder(res).Encode(message)
		}
	} else {
		if VerifyTokenPointer(tokenPointer) {
			if tokenPointer.Valid {
				json.NewEncoder(res).Encode(message)
			} else {
				if message == "Success" {
					user, _ := db.GetUserData(email)
					token, _ := GenerateToken(user, time.Hour)
					json.NewEncoder(res).Encode(token)
				} else {
					json.NewEncoder(res).Encode(message)
				}
			}
		} else {
			json.NewEncoder(res).Encode("Token not found")
		}
	}
}

// Logic for handling user login via Google OAuth authentication, including fetching user information using the access token.
func TokenGoogleLoginLogic(res http.ResponseWriter, req *http.Request, accessToken string) *oauth2v2.Userinfo {
	if accessToken == "" {
		http.Error(res, "Unauthorised", http.StatusBadRequest)
		return nil
	}
	config := oauth2.Config{}
	token := &oauth2.Token{AccessToken: accessToken}
	ctx := context.Background()
	client := config.Client(ctx, token)

	service, err := oauth2v2.New(client)
	if err != nil {
		log.Fatalf("Error creating OAuth2 service: %v", err)
	}

	info, err := service.Userinfo.Get().Do()
	if err != nil {
		log.Fatalf("Error getting user info: %v", err)
	}

	return info
}
