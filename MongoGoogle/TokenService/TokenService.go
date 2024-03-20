package TokenService

import (
	model "MongoGoogle/Model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	applicationService "MongoGoogle/ApplicationService"
	db "MongoGoogle/Repository"

	oauth2v2 "google.golang.org/api/oauth2/v2"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

var jwtKey = []byte("tajna_lozinka")

// Korisnik predstavlja strukturu korisnika

// GenerateToken generira JWT token na temelju korisničkih podataka
func GenerateToken(user model.ApplicationUser, exp time.Duration) (string, error) {
	//tokenTTL := 1 * time.Minute // Token vredi 1 sat
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":                user.ID,
		"email":             user.Email,
		"firstName":         user.FirstName,
		"lastName":          user.LastName,
		"phone":             user.Phone,
		"dateOfBirth":       user.DateOfBirth,
		"username":          user.Username,
		"password":          user.Password,
		"role":              user.Role,
		"verified":          user.Verified,
		"applicationMethod": user.ApplicationMethod,
		"exp":               time.Now().Add(exp).Unix(),
	})

	// Potpisivanje tokena
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

func ParseTokenString(tokenString string) (*jwt.Token, error) {
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

func VerifyTokenPointer(token *jwt.Token) bool {
	if token == nil {
		return false
	} else {
		return true
	}
}

// Extracts user information from the token in the request header.
func ExtractUserFromToken(tokenString string) (model.ApplicationUser, *jwt.Token, error) {
	var usererr model.ApplicationUser
	var errtoken *jwt.Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return usererr, errtoken, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	verified, _ := claims["verified"].(bool)

	id, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	user := model.ApplicationUser{
		ID:          id,
		Email:       claims["email"].(string),
		FirstName:   claims["firstName"].(string),
		LastName:    claims["lastName"].(string),
		Phone:       claims["phone"].(string),
		DateOfBirth: claims["dateOfBirth"].(string),
		Username:    claims["username"].(string),
		Password:    claims["password"].(string),
		Role:        claims["role"].(string),
		Verified:    verified,
	}
	return user, token, nil
}

func TokenAppLoginLogic(res http.ResponseWriter, req *http.Request, authHeader string, email string, password string) {

	tokenString := SplitTokenHeder(authHeader)
	tokenPointer, _ := ParseTokenString(tokenString)
	message := applicationService.ApplicationLogin(email, password)

	if tokenString == "" {
		if message == "Success" {
			user, _ := db.GetUserData(email)
			token, _ := GenerateToken(user, time.Hour)
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode("This is your bearer token for login: " + token)
		} else {
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode("Unauthorised")
		}
	} else {
		if VerifyTokenPointer(tokenPointer) {
			if tokenPointer.Valid {
				res.Header().Set("Content-Type", "application/json")
				json.NewEncoder(res).Encode(message)
			} else {
				if message == "Success" {
					user, _ := db.GetUserData(email)
					token, _ := GenerateToken(user, time.Hour)
					res.Header().Set("Content-Type", "application/json")
					json.NewEncoder(res).Encode("This is your bearer token for login: " + token)
				} else {
					res.Header().Set("Content-Type", "application/json")
					json.NewEncoder(res).Encode(message)
				}
			}
		} else {
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode("Token not found")
		}
	}
}

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
