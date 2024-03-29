package Controller

import (
	service "MongoGoogle/Service"
	"net/http"

	"github.com/gorilla/mux"
)

type CompanyController struct {
	Router *mux.Router
}

func NewCompanyController() *CompanyController {
	return &CompanyController{
		Router: mux.NewRouter(),
	}
}
func (cc *CompanyController) RegisterRoutes() {
	cc.Router.HandleFunc("/company/registerCompany", cc.RegisterCompany)
	cc.Router.HandleFunc("/company/deleteCompany", cc.DeleteCompany)
}

func (cc *CompanyController) RegisterCompany(res http.ResponseWriter, req *http.Request) {
	service.CreateComapny(res, req)
}
func (cc *CompanyController) DeleteCompany(res http.ResponseWriter, req *http.Request) {
	service.CreateComapny(res, req)
}
