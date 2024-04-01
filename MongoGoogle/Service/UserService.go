package Service

import (
	model "MongoGoogle/Model"
	repo "MongoGoogle/Repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserData(mail string) model.ApplicationUser {
	user, _ := repo.GetUserData(mail)
	return user
}

func VerifyUser(email string) bool {
	success := repo.VerifyUser(email)
	return success
}

func CheckUserRole(user model.ApplicationUser, role string) (bool, primitive.ObjectID) {
	userComapnies := user.Companies
	for _, value := range userComapnies {
		if value.Role == role {
			return true, value.CompanyID
		}
	}
	return false, primitive.ObjectID{}
}
