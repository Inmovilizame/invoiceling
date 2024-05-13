package service

import (
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type FreelancerRepo interface {
	Get() *model.Freelancer
}

type Freelancer struct {
	repo FreelancerRepo
}

func NewFreelancer(repo repository.CfgFreelancer) *Freelancer {
	return &Freelancer{
		repo: repo,
	}
}

func (fs *Freelancer) Get() *model.Freelancer {
	return fs.repo.Get()
}
