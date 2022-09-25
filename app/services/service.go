package services

import (
	"encoding/json"
	"fmt"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"io/ioutil"
	"net/http"
)

type Service struct {
	Repo      usersRepository
	ConvertTo string
}

type usersRepository interface {
	IsUserExisting(id int) bool
	GetUserBalance(id int) (float64, error)
	DepositMoney(user users.User) error
	WithdrawMoney(user users.User) error
	TransferMoney(usersTransfer users.TransferMoney) error
	CreateUser(id int) error
}

func New(repository usersRepository) *Service {
	serv := Service{
		Repo:      repository,
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
		defer res.Body.Close()
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
	err := service.Repo.DepositMoney(user)
	return err
}

func (service *Service) WithdrawMoney(user users.User) error {
	err := service.Repo.WithdrawMoney(user)
	return err
}

func (service *Service) TransferMoney(usersTransfer users.TransferMoney) error {
	err := service.Repo.TransferMoney(usersTransfer)
	return err
}

func (service *Service) IsUserExisting(id int) bool {
	isExist := service.Repo.IsUserExisting(id)
	return isExist
}

func (service *Service) CreateUser(id int) error {
	err := service.Repo.CreateUser(id)
	return err
}
