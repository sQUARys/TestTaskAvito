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
	sync.Mutex
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
		status, err := w.Write([]byte("Error strconv in controller level : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	if !ctr.Service.IsUserExisting(idInt) { // если в бд нет пользователя с таким id выводим ошибку
		status, err := w.Write([]byte("This user doesn't exist, choose another existing user. "))
		ErrorHandler(w, status, err)
		return
	}

	balance, err := ctr.Service.GetUserBalance(idInt)

	if err != nil {
		status, err := w.Write([]byte("Error balance in controller level : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		status, err := w.Write([]byte("Error json in controller level : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	_, err = w.Write(balanceJSON)
	if err != nil {
		status, err := w.Write([]byte("Error json in controller level : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}
}

func (ctr *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		status, err := w.Write([]byte("Error strconv in controller level : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	err = ctr.Service.CreateUser(idInt)
	if err != nil {
		status, err := w.Write([]byte("Error in creating user :  " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}
}

func (ctr *Controller) DepositMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status, err := w.Write([]byte("Error in deposit money :" + err.Error()))
		ErrorHandler(w, status, err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		status, err := w.Write([]byte("Error in deposit money :" + err.Error()))
		ErrorHandler(w, status, err)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		status, err := w.Write([]byte("This user doesn't exist, choose another existing user. " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	if user.UpdateValue <= 0 {
		status, err := w.Write([]byte("Deposit value can't be negative or nought. " + err.Error()))
		ErrorHandler(w, status, err)
	}

	ctr.Lock()
	defer ctr.Unlock() // для того чтобы при нескольких одновременно отправленных запросов на внесение, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.DepositMoney(user)

	if err != nil {
		status, err := w.Write([]byte("Error in deposit money : " + err.Error()))
		ErrorHandler(w, status, err)
	}
}

func (ctr *Controller) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status, err := w.Write([]byte("Error in withdraw money : " + err.Error()))
		ErrorHandler(w, status, err)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		status, err := w.Write([]byte("Error in withdraw money : " + err.Error()))
		ErrorHandler(w, status, err)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		status, err := w.Write([]byte("This user doesn't exist, choose another existing user. " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	if user.UpdateValue <= 0 {
		status, err := w.Write([]byte("Withdraw value can't be negative or nought. " + err.Error()))
		ErrorHandler(w, status, err)
	}

	ctr.Lock()
	defer ctr.Unlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.WithdrawMoney(user)
	if err != nil {
		status, err := w.Write([]byte("Error in withdraw money : " + err.Error()))
		ErrorHandler(w, status, err)
	}
}

func (ctr *Controller) TransferMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status, err := w.Write([]byte("Error in transfer money : " + err.Error()))
		ErrorHandler(w, status, err)
	}

	var usersTransfer users.TransferMoney

	err = json.Unmarshal(body, &usersTransfer)
	if err != nil {
		status, err := w.Write([]byte("Error in transfer money : " + err.Error()))
		ErrorHandler(w, status, err)
	}

	if !ctr.Service.IsUserExisting(usersTransfer.IdOfSenderUser) || !ctr.Service.IsUserExisting(usersTransfer.IdOfRecipientUser) { // если в бд нет пользователя с таким id выводим ошибку
		status, err := w.Write([]byte("This user doesn't exist, choose another existing user. " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	if usersTransfer.SendingAmount <= 0 {
		status, err := w.Write([]byte("Sending amount value can't be negative or nought. " + err.Error()))
		ErrorHandler(w, status, err)
	}

	ctr.Lock()
	defer ctr.Unlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.TransferMoney(usersTransfer)
	if err != nil {
		status, err := w.Write([]byte("Error in transfer money : " + err.Error()))
		ErrorHandler(w, status, err)
	}

}

func (ctr *Controller) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		status, err := w.Write([]byte("Error strconv in controller level : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	transactions, err := ctr.Service.GetUserTransactions(idInt)

	if err != nil {
		status, err := w.Write([]byte("Error in get transactions : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}

	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		status, err := w.Write([]byte("Error in json transactions : " + err.Error()))
		ErrorHandler(w, status, err)
		return
	}
	w.Write(transactionsJSON)
}

func ErrorHandler(w http.ResponseWriter, status int, err error) {
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
}
