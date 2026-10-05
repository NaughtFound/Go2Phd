package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"naughtfound.github.io/go2phd/models"
	"naughtfound.github.io/go2phd/repositories"
)

type PositionService struct {
	uniRepo      *repositories.UniversityRepository
	clientID     string
	clientSecret string
}

func NewPositionService(uniRepo *repositories.UniversityRepository, clientID, clientSecret string) *PositionService {
	return &PositionService{
		uniRepo:      uniRepo,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *PositionService) FetchAllPositions(refreshToken, spreadsheetID, sheetName string) ([]models.Position, error) {

	sheetsRepo, err := repositories.NewSheetsRepository(s.clientID, s.clientSecret, refreshToken, spreadsheetID)
	if err != nil {
		return nil, err
	}

	rows, err := sheetsRepo.ReadAllRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) <= 1 {
		return []models.Position{}, nil
	}

	var positions []models.Position

	for i, row := range rows[1:] {
		pos := parseRowToPosition(i+1, row)

		if pos.University.Name != "" {
			if uni, err := s.uniRepo.GetByName(pos.University.Name); err == nil {
				pos.University = uni
			}
		}

		positions = append(positions, pos)
	}

	return positions, nil
}

func parseRowToPosition(id int, row []any) models.Position {
	pos := models.Position{
		ID: id,
	}

	getCol := func(idx int) string {
		if idx < len(row) && row[idx] != nil {
			return strings.TrimSpace(fmt.Sprintf("%v", row[idx]))
		}
		return ""
	}

	// Column Mapping (assuming sheet column order):
	// Col 0: Title
	// Col 1: Priority
	// Col 2: Tags (comma-separated: "AI, Go, Microservices")
	// Col 3: Status
	// Col 4: University Name
	// Col 5: Due date (e.g. "2026-11-15" or "2026-11-15T00:00:00Z")
	// Col 6: RollingBased ("true"/"false" or "TRUE"/"FALSE")
	// Col 7: Links (comma-separated or single URL)

	pos.Title = getCol(0)
	pos.Priority = models.Priority(getCol(1))

	rawTags := getCol(2)
	if rawTags != "" {
		tagList := strings.SplitSeq(rawTags, ",")
		for t := range tagList {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				pos.Tags = append(pos.Tags, trimmed)
			}
		}
	}

	pos.Status = models.Status(getCol(3))

	uniName := getCol(4)
	if uniName != "" {
		pos.University = models.University{Name: uniName}
	}

	rawDueDate := getCol(5)
	if rawDueDate != "" {
		if parsedDate, err := parseDate(rawDueDate); err == nil {
			pos.DueDate = &parsedDate
		}
	}

	rawRolling := strings.ToLower(getCol(6))
	if b, err := strconv.ParseBool(rawRolling); err == nil {
		pos.RollingBased = b
	} else if rawRolling == "yes" || rawRolling == "y" {
		pos.RollingBased = true
	}

	rawLinks := getCol(7)
	if rawLinks != "" {
		linkList := strings.SplitSeq(rawLinks, ",")
		for l := range linkList {
			if trimmed := strings.TrimSpace(l); trimmed != "" {
				pos.Links = append(pos.Links, trimmed)
			}
		}
	}

	return pos
}

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		"01/02/2006",
		"02/01/2006",
		"2006/01/02",
	}

	for _, fmtStr := range formats {
		if t, err := time.Parse(fmtStr, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}
