package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sQUARys/TestTaskAvito/app/services"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Controller struct {
	Service services.Service
	mut     sync.Mutex
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

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("Error strconv in controller level : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ctr.Service.IsUserExisting(idInt) { // если в бд нет пользователя с таким id выводим ошибку
		log.Println("This user doesn't exist, choose another existing user. ")
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

	status, err := w.Write(balanceJSON)
	if err != nil {
		log.Println("Error json in controller level : ", err)
		w.WriteHeader(status)
		return
	}

}

func (ctr *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("Error strconv in controller level : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ctr.Service.CreateUser(idInt)
	if err != nil {
		log.Println("Error in creating user : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctr *Controller) DepositMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in deposit money : ", err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error in deposit money : ", err)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		log.Println("This user doesn't exist, choose another existing user. ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.UpdateValue <= 0 {
		log.Println("Deposit value can't be negative or nought. ")
	}

	ctr.mut.Lock()
	defer ctr.mut.Unlock() // для того чтобы при нескольких одновременно отправленных запросов на внесение, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.DepositMoney(user)

	if err != nil {
		log.Println("Error in deposit money : ", err)
	}
}

func (ctr *Controller) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in withdraw money : ", err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("Error in withdraw money : ", err)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		log.Println("This user doesn't exist, choose another existing user. ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.UpdateValue <= 0 {
		log.Println("Withdraw value can't be negative or nought. ")
	}

	ctr.mut.Lock()
	defer ctr.mut.Unlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.WithdrawMoney(user)
	if err != nil {
		log.Println("Error in withdraw money : ", err)
	}
}

func (ctr *Controller) TransferMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in transfer money : ", err)
	}

	var usersTransfer users.TransferMoney

	err = json.Unmarshal(body, &usersTransfer)
	if err != nil {
		log.Println("Error in transfer money : ", err)
	}

	if !ctr.Service.IsUserExisting(usersTransfer.IdOfSenderUser) || !ctr.Service.IsUserExisting(usersTransfer.IdOfRecipientUser) { // если в бд нет пользователя с таким id выводим ошибку
		log.Println("This user doesn't exist, choose another existing user. ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if usersTransfer.SendingAmount <= 0 {
		log.Println("Sending amount value can't be negative or nought. ")
	}

	ctr.mut.Lock()
	defer ctr.mut.Unlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.TransferMoney(usersTransfer)
	if err != nil {
		log.Println("Error in withdraw money : ", err)
	}

}
