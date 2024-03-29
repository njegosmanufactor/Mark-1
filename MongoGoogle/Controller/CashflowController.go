package Controller

import (
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

}
