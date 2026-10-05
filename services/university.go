package services

import (
	"naughtfound.github.io/go2phd/models"
	"naughtfound.github.io/go2phd/repositories"
)

type UniversityService struct {
	repo *repositories.UniversityRepository
}

func NewUniversityService(repo *repositories.UniversityRepository) *UniversityService {
	return &UniversityService{repo: repo}
}

func (s *UniversityService) GetUniversityInfo(name string) (models.University, error) {
	return s.repo.GetByName(name)
}
