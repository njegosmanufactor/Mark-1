package Service

import (
	model "MongoGoogle/Model"
	repo "MongoGoogle/Repository"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeletePendingRequest(id string) {
	repo.DeletePandingRequrst(id)
}

func FindCodeRequestByHex(id string, res http.ResponseWriter) (model.PasswordLessRequest, bool) {
	result, found := repo.FindCodeRequestByHex(id, res)
	return result, found
}

func CheckForUnregInvites(email string, res http.ResponseWriter) (primitive.ObjectID, bool) {
	id, found := repo.FindUnregInvite(email, res)
	return id, found
}
