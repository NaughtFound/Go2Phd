package services

import (
	"naughtfound.github.io/go2phd/models"
	"naughtfound.github.io/go2phd/repositories"
)

type PositionService struct {
	sheetsRepo *repositories.SheetsRepository
	uniRepo    *repositories.UniversityRepository
}

func NewPositionService(sheetsRepo *repositories.SheetsRepository, uniRepo *repositories.UniversityRepository) *PositionService {
	return &PositionService{
		sheetsRepo: sheetsRepo,
		uniRepo:    uniRepo,
	}
}

func (s *PositionService) FetchAllPositions(sheetName string) ([]models.Position, error) {
	positions, err := s.sheetsRepo.GetAllPositions(sheetName)
	if err != nil {
		return nil, err
	}

	for i := range positions {
		if positions[i].University.Name != "" {
			uni, err := s.uniRepo.GetByName(positions[i].University.Name)
			if err == nil {
				positions[i].University = uni
			}
		}
	}

	return positions, nil
}
