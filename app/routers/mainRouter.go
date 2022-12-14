package routers

import (
	"github.com/gorilla/mux"
	controller "github.com/sQUARys/TestTaskAvito/app/controllers"
)

type Router struct {
	Router     *mux.Router
	Controller controller.Controller
}

func New(controller *controller.Controller) *Router {
	r := mux.NewRouter()
	return &Router{
		Controller: *controller,
		Router:     r,
	}
}

func (r *Router) SetRoutes() {
	r.Router.HandleFunc("/get-balance/{id}&{currency}", r.Controller.GetUserBalance).Methods("Get")
	r.Router.HandleFunc("/create-user/{id}", r.Controller.CreateUser).Methods("Get")
	r.Router.HandleFunc("/get-transactions/{id}", r.Controller.GetUserTransactions).Methods("Get")

	r.Router.HandleFunc("/deposit-money/", r.Controller.DepositMoney).Methods("Post")
	r.Router.HandleFunc("/withdraw-money/", r.Controller.WithdrawMoney).Methods("Post")
	r.Router.HandleFunc("/transfer-money/", r.Controller.TransferMoney).Methods("Post")

}
