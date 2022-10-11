package services

import (
	"encoding/json"
	"fmt"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	Repo      usersRepository
	Cache     transactionsCache
	ConvertTo string
	sync.Mutex
}

type transactionsCache interface {
	GetUserTransactions(id int) ([]users.Transaction, error)
	AddTransaction(userId int, description string, date string) error
}
type usersRepository interface {
	IsUserExisting(id int) bool
	GetUserBalance(id int) (float64, error)
	DepositMoney(user users.User) error
	WithdrawMoney(user users.User) error
	TransferMoney(usersTransfer users.TransferMoney) error
	CreateUser(id int) error
}

func New(repository usersRepository, cache transactionsCache) *Service {
	serv := Service{
		Repo:      repository,
		Cache:     cache,
		ConvertTo: "",
	}
	return &serv
}

func (service *Service) Convert(from string, to string, amount float64) (float64, error) {
	format := "https://api.apilayer.com/exchangerates_data/convert?to=%s&from=%s&amount=%f"

	url := fmt.Sprintf(format, to, from, amount)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("apikey", "OW7XhQPjG87aX7vEg5bz0glUbL2BgMeX")

	res, err := client.Do(req)
	if res.Body != nil {
		defer func() {
			err := res.Body.Close()
			if err != nil {
				log.Println(err)
			}
		}()
	}

	jsonMap := make(map[string]interface{})
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return 0, err
	}

	return jsonMap["result"].(float64), nil
}

func (service *Service) GetUserBalance(id int) (float64, error) {
	service.Lock()
	defer service.Unlock()

	balance, err := service.Repo.GetUserBalance(id)
	if err != nil {
		return -1, err
	}
	if service.ConvertTo != "" {
		balance, err = service.Convert("RUB", service.ConvertTo, balance)
		if err != nil {
			return -1, err
		}
	}
	return balance, err
}

func (service *Service) DepositMoney(user users.User) error {
	service.Lock()
	defer service.Unlock()
	err := service.Repo.DepositMoney(user)
	if err == nil { // если все успешно внеслось
		err = service.Cache.AddTransaction(user.Id, fmt.Sprintf("Внесение на сумму: %2f", user.UpdateValue),
			time.Now().Format("01-02-2006 15:04:05"))
		if err != nil {
			return err
		}
	}
	return err
}

func (service *Service) WithdrawMoney(user users.User) error {
	service.Lock()
	defer service.Unlock()

	err := service.Repo.WithdrawMoney(user)
	if err == nil { // если все успешно внеслось
		err = service.Cache.AddTransaction(user.Id, fmt.Sprintf("Снятие на сумму: %2f", user.UpdateValue), time.Now().Format("01-02-2006 15:04:05"))
		if err != nil {
			return err
		}
	}
	return err
}

func (service *Service) TransferMoney(usersTransfer users.TransferMoney) error {
	service.Lock()
	defer service.Unlock()

	err := service.Repo.TransferMoney(usersTransfer)
	if err == nil { // если все успешно внеслось
		err = service.AddTransaction(usersTransfer.IdOfSenderUser, fmt.Sprintf("Перевод на сумму: %2f", usersTransfer.SendingAmount), time.Now().Format("01-02-2006 15:04:05"))
		if err != nil {
			return err
		}
		err = service.AddTransaction(usersTransfer.IdOfRecipientUser, fmt.Sprintf("Получение перевода на сумму: %2f", usersTransfer.SendingAmount), time.Now().Format("01-02-2006 15:04:05"))
		if err != nil {
			return err
		}
	}
	return err
}

func (service *Service) IsUserExisting(id int) bool {
	service.Lock()
	defer service.Unlock()

	isExist := service.Repo.IsUserExisting(id)
	return isExist
}

func (service *Service) CreateUser(id int) error {
	service.Lock()
	defer service.Unlock()

	err := service.Repo.CreateUser(id)
	return err
}

func (service *Service) AddTransaction(id int, description string, date string) error {
	err := service.Cache.AddTransaction(id, description, date)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) GetUserTransactions(id int) ([]users.Transaction, error) {
	service.Lock()
	defer service.Unlock()

	transactions, err := service.Cache.GetUserTransactions(id)
	return transactions, err
}
