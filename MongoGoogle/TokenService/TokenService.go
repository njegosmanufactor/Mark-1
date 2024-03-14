package TokenService

import (
	model "MongoGoogle/Model"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = []byte("tajna_lozinka")

// Korisnik predstavlja strukturu korisnika

// GenerateToken generira JWT token na temelju korisničkih podataka
func GenerateToken(user model.ApplicationUser) (string, error) {
	tokenTTL := 1 * time.Hour // Token vrijedi 1 sat, možete prilagoditi ovo prema potrebama
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          user.ID,
		"email":       user.Email,
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
		"phone":       user.Phone,
		"dateOfBirth": user.DateOfBirth,
		"username":    user.Username,
		"password":    user.Password,
		"company":     user.Company,
		"role":        user.Role,
		"verified":    user.Verified,
		"exp":         time.Now().Add(tokenTTL).Unix(),
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
		return "Invalid or missing Bearer token"
	}
	tokenString := parts[1]
	return tokenString
}

// Extracts user information from the token in the request header and sets the user as authorized in the database.
func ExtractUserFromToken(tokenString string) (model.ApplicationUser, *jwt.Token, error) {
	var usererr model.ApplicationUser
	var errtoken *jwt.Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Provjera metode potpisivanja
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil // jwtKey je globalna varijabla definirana u prethodnom primjeru
	})
	if err != nil {
		return usererr, errtoken, fmt.Errorf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Errorf("Invalid token claims")
	}

	verified, ok := claims["verified"].(bool)
	if !ok {
		fmt.Errorf("Invalid or missing 'verified' claim")
	}

	id, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	// Kreiranje ApplicationUser objekta
	user := model.ApplicationUser{
		ID:          id,
		Email:       claims["email"].(string),
		FirstName:   claims["firstName"].(string),
		LastName:    claims["lastName"].(string),
		Phone:       claims["phone"].(string),
		DateOfBirth: claims["dateOfBirth"].(string),
		Username:    claims["username"].(string),
		Password:    claims["password"].(string),
		Company:     claims["company"].(string),
		Role:        claims["role"].(string),
		Verified:    verified,
	}

	fmt.Println(user)
	return user, token, nil
}
