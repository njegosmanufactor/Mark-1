package Controller

import (
	"encoding/json"
	"net/http"

	service "MongoGoogle/Service"

	"github.com/gorilla/mux"
)

type UserController struct {
	Router *mux.Router
}

func NewUserController() *UserController {
	return &UserController{
		Router: mux.NewRouter(),
	}
}

func (uc *UserController) RegisterRoutes() {
	uc.Router.HandleFunc("/users/sendInvitation", uc.SendInvitation)
	uc.Router.HandleFunc("/users/inviteConfirmation/{id}", uc.SendInvitationConfirmation)
	uc.Router.HandleFunc("/users/forgotPassword", uc.ChangePassword)
	uc.Router.HandleFunc("/users/forgotPassword/callback/{transferId}", uc.ChangePassword)
	uc.Router.HandleFunc("/users/transferOwnership", uc.TransferOwnership)
	uc.Router.HandleFunc("/users/transferOwnership/feedback/{transferId}", uc.TransferOwnershipCallback)

}

func (uc *UserController) SendInvitation(res http.ResponseWriter, req *http.Request) {
	user, tokenpointer := service.GetUserAndPointerFromToken(res, req)
	if tokenpointer != nil && tokenpointer.Valid {
		service.SendInvitation(res, req, user)
	} else {
		json.NewEncoder(res).Encode("Session timed out or terminated")
	}
}
func (uc *UserController) SendInvitationConfirmation(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	transactionId := vars["id"]
	service.IncludeUserInCompany(transactionId, res)
}

func (uc *UserController) ChangePassword(res http.ResponseWriter, req *http.Request) {
	_, tokenpointer := service.GetUserAndPointerFromToken(res, req)
	if tokenpointer != nil && tokenpointer.Valid {
		service.PasswordChange(res, req)
	} else {
		json.NewEncoder(res).Encode("Session timed out or terminated")
	}
}
func (uc *UserController) ChangePasswordHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	transferId := vars["transferId"]
	service.FinaliseForgottenPasswordUpdate(transferId, res, req)
}
func (uc *UserController) TransferOwnership(res http.ResponseWriter, req *http.Request) {
	user, tokenpointer := service.GetUserAndPointerFromToken(res, req)

	if tokenpointer != nil && tokenpointer.Valid {
		service.TransferOwnership(res, req, user)
	} else {
		json.NewEncoder(res).Encode("Session timed out or terminated")
	}
}
func (uc *UserController) TransferOwnershipCallback(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	transferId := vars["transferId"]
	service.FinaliseOwnershipTransfer(transferId, res)
}
