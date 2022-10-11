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
	sync.RWMutex
}

func New(service *services.Service) *Controller {
	return &Controller{
		Service: *service,
	}
}

func (ctr *Controller) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	currency := vars["currency"]
	ctr.Service.ConvertTo = currency
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	if !ctr.Service.IsUserExisting(idInt) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err)
		return
	}

	balance, err := ctr.Service.GetUserBalance(idInt)

	if err != nil {
		ErrorHandler(w, err)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	_, err = w.Write(balanceJSON)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
}

func (ctr *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	if idInt < 0 {
		ErrorHandler(w, err)
		return
	}
	err = ctr.Service.CreateUser(idInt)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
}

func (ctr *Controller) DepositMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		ErrorHandler(w, err)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err)
	}

	if user.UpdateValue <= 0 {
		ErrorHandler(w, err)
	}

	ctr.RLock()
	defer ctr.RUnlock() // для того чтобы при нескольких одновременно отправленных запросов на внесение, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.DepositMoney(user)

	if err != nil {
		ErrorHandler(w, err)
	}
}

func (ctr *Controller) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		ErrorHandler(w, err)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err)
		return
	}

	if user.UpdateValue <= 0 {
		ErrorHandler(w, err)
	}

	ctr.RLock()
	defer ctr.RUnlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.WithdrawMoney(user)
	if err != nil {
		ErrorHandler(w, err)
	}
}

func (ctr *Controller) TransferMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err)
	}

	var usersTransfer users.TransferMoney

	err = json.Unmarshal(body, &usersTransfer)
	if err != nil {
		ErrorHandler(w, err)
	}

	if !ctr.Service.IsUserExisting(usersTransfer.IdOfSenderUser) || !ctr.Service.IsUserExisting(usersTransfer.IdOfRecipientUser) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err)
		return
	}

	if usersTransfer.SendingAmount <= 0 {
		ErrorHandler(w, err)
	}

	ctr.RLock()
	defer ctr.RUnlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.TransferMoney(usersTransfer)
	if err != nil {
		ErrorHandler(w, err)
	}

}

func (ctr *Controller) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	transactions, err := ctr.Service.GetUserTransactions(idInt)

	if err != nil {
		ErrorHandler(w, err)
		return
	}

	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	w.Write(transactionsJSON)
}

func ErrorHandler(w http.ResponseWriter, err error) {
	_, writeError := w.Write([]byte(err.Error()))
	if writeError != nil {
		log.Println(writeError)
	}
	w.WriteHeader(status)
}
