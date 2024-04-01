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
	cfc.Router.HandleFunc("/cashFlow/CreateInflow", cfc.CreateInflow)
}

func (cc *CashflowController) CreateCashFlowForUser(res http.ResponseWriter, req *http.Request) {
	service.CreateCashFlowUser(res, req)
}

func (cc *CashflowController) CreateInflow(res http.ResponseWriter, req *http.Request) {
	service.CreateInflow(res, req)
}
