package Controller

import (
	service "MongoGoogle/Service"
	"net/http"

	"github.com/gorilla/mux"
)

type CashflowController struct {
	Router *mux.Router
}

func NewCashflowController() *CashflowController {
	return &CashflowController{
		Router: mux.NewRouter(),
	}
}
func (cfc *CashflowController) RegisterRoutes() {
	cfc.Router.HandleFunc("/cashFlow/CreateCashFlowUser", cfc.CreateCashFlowForUser)
}

func (cc *CashflowController) CreateCashFlowForUser(res http.ResponseWriter, req *http.Request) {
	service.CreateCashFlowUser(res, req)
}
