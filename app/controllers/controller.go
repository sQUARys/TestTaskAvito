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

const (
	ok             = http.StatusOK
	noContent      = http.StatusNoContent
	serverInternal = http.StatusInternalServerError
	badRequest     = http.StatusBadRequest
	notFound       = http.StatusNotFound
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
		ErrorHandler(w, err, serverInternal)
		return
	}

	if !ctr.Service.IsUserExisting(idInt) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err, noContent)
		return
	}

	balance, err := ctr.Service.GetUserBalance(idInt)

	if err != nil {
		ErrorHandler(w, err, noContent)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	w.WriteHeader(ok)
	_, err = w.Write(balanceJSON)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}
}

func (ctr *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}
	if idInt < 0 {
		ErrorHandler(w, err, badRequest)
		return
	}
	err = ctr.Service.CreateUser(idInt)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}
	w.WriteHeader(ok)
}

func (ctr *Controller) DepositMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err, notFound)
	}

	if user.UpdateValue <= 0 {
		ErrorHandler(w, err, badRequest)
	}

	ctr.RLock()
	defer ctr.RUnlock() // для того чтобы при нескольких одновременно отправленных запросов на внесение, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.DepositMoney(user)

	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}
	w.WriteHeader(ok)
}

func (ctr *Controller) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	var user users.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	if !ctr.Service.IsUserExisting(user.Id) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err, notFound)
		return
	}

	if user.UpdateValue <= 0 {
		ErrorHandler(w, err, badRequest)
	}

	ctr.RLock()
	defer ctr.RUnlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.WithdrawMoney(user)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}
	w.WriteHeader(ok)
}

func (ctr *Controller) TransferMoney(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	var usersTransfer users.TransferMoney

	err = json.Unmarshal(body, &usersTransfer)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	if !ctr.Service.IsUserExisting(usersTransfer.IdOfSenderUser) || !ctr.Service.IsUserExisting(usersTransfer.IdOfRecipientUser) { // если в бд нет пользователя с таким id выводим ошибку
		ErrorHandler(w, err, notFound)
		return
	}

	if usersTransfer.SendingAmount <= 0 {
		ErrorHandler(w, err, badRequest)
	}

	ctr.RLock()
	defer ctr.RUnlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)

	err = ctr.Service.TransferMoney(usersTransfer)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}
	w.WriteHeader(ok)
}

func (ctr *Controller) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	transactions, err := ctr.Service.GetUserTransactions(idInt)

	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}
	w.WriteHeader(ok)
	w.Write(transactionsJSON)
}

func ErrorHandler(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_, writeError := w.Write([]byte(err.Error()))
	if writeError != nil {
		log.Println(writeError)
	}
}
