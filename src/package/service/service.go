package service

import (
	"src/models"
	"src/package/repository"
)

type Authorization interface {
	CreateUser(user models.User) (*models.User, error)
	GenerateToken(user models.User) (string, error)
	ParseToken(token string) (int, error)
}

type Surveys interface{}

type Service struct {
	Authorization
	Surveys
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(&repos.Authorization),
	}
}
