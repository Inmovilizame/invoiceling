package service

import (
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type FreelancerRepo interface {
	GetFreelancer() *model.Freelancer
}

type Freelancer struct {
	repo FreelancerRepo
}

func NewFreelancer(repo repository.CfgRepo) *Freelancer {
	return &Freelancer{
		repo: repo,
	}
}

func (fs *Freelancer) GetFreelancer() *model.Freelancer {
	return fs.repo.GetFreelancer()
}
