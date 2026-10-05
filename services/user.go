package services

import (
	"naughtfound.github.io/go2phd/models"
	"naughtfound.github.io/go2phd/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (models.User, error) {
	return s.repo.FindByID(id)
}
