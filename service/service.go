package service

import "new/repositories/orders_repo"

type Service struct {
	Repositories Repositories
}

type Repositories struct {
	Orders *orders_repo.Repo
}

func New(repos Repositories) *Service {
	return &Service{
		Repositories: repos,
	}
}
