package service

import "new/repositories/orders_repo"

type Sevice struct {
	Repositories Repositories
}

type Repositories struct {
	Orders *orders_repo.Repo
}

func New(repos Repositories) *Sevice {
	return &Sevice{
		Repositories: repos,
	}
}
