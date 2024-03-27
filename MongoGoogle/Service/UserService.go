package Service

import (
	model "MongoGoogle/Model"
	repo "MongoGoogle/Repository"
)

func GetUserData(mail string) model.ApplicationUser {
	user, _ := repo.GetUserData(mail)
	return user
}

func VerifyUser(email string) bool {
	success := repo.VerifyUser(email)
	return success
}
