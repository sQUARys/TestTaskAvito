package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sQUARys/TestTaskAvito/app/services"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Controller struct {
	Service services.Service
}

func New(service *services.Service) *Controller {
	return &Controller{
		Service: *service,
	}
}

func (ctr *Controller) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	idString := vars["id"]
	fmt.Println("ID : ", idString)

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("Error strconv in controller level : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := ctr.Service.GetUserBalance(idInt)

	if err != nil {
		log.Println("Error balance in controller level : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		log.Println("Error json in controller level : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(balanceJSON)

}

func (ctr *Controller) DepositMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error in deposit money : ", err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Error in deposit money : ", err)
	}

	if user.UpdateValue < 0 {
		fmt.Println("Deposit value can't be negative. ")
	}

	err = ctr.Service.DepositMoney(user)
	if err != nil {
		fmt.Println("Error in deposit money : ", err)
	}
}

func (ctr *Controller) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error in withdraw money : ", err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Error in withdraw money : ", err)
	}

	err = ctr.Service.WithdrawMoney(user)
	if err != nil {
		fmt.Println("Error in withdraw money : ", err)
	}
}
