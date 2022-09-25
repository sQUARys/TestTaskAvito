package services

import (
	"github.com/sQUARys/TestTaskAvito/app/users"
)

type Service struct {
	Repo usersRepository
}

type usersRepository interface {
	GetUserBalance(id int) (int, error)
	DepositMoney(user users.User) error
	WithdrawMoney(user users.User) error
	TransferMoney(usersTransfer users.TransferMoney) error
	isUserExisting(id int) bool
}

func New(repository usersRepository) *Service {
	serv := Service{
		Repo: repository,
	}
	return &serv
}

func (service *Service) GetUserBalance(id int) (int, error) {
	balance, err := service.Repo.GetUserBalance(id)
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

func (service *Service) isUserExisting(id int) bool {
	isExist := service.Repo.isUserExisting(id)
	return isExist
}
