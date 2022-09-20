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

func (service *Service) DepositMoney(user users.User) {
	if user.Deposit < 0 {
		return
	}
	service.Repo.DepositMoney(user)
}
