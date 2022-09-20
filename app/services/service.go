package services

type Service struct {
	Repo usersRepository
}

type usersRepository interface {
	GetUserBalance(id int) (int, error)
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
