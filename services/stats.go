package services

import (
	"strings"

	"naughtfound.github.io/go2phd/models"
)

type StatsService struct {
	positionService *PositionService
}

func NewStatsService(posService *PositionService) *StatsService {
	return &StatsService{positionService: posService}
}

func CalculateOverallStats(positions []models.Position) models.OverallStats {
	statusCounts := map[models.Status]int{
		models.StatusNotStarted:  0,
		models.StatusInProgress:  0,
		models.StatusApplied:     0,
		models.StatusEmailed:     0,
		models.StatusInterview:   0,
		models.StatusRejected:    0,
		models.StatusShortlisted: 0,
		models.StatusAccepted:    0,
		models.StatusCanceled:    0,
		models.StatusNoAnswer:    0,
	}

	priorityCounts := map[models.Priority]int{
		models.PriorityHigh:   0,
		models.PriorityMedium: 0,
		models.PriorityLow:    0,
	}

	tagsCounts := make(map[string]int)
	rollingCount := 0

	for _, pos := range positions {
		if pos.Status != "" {
			statusCounts[pos.Status]++
		}
		if pos.Priority != "" {
			priorityCounts[pos.Priority]++
		}
		if pos.RollingBased {
			rollingCount++
		}
		for _, tag := range pos.Tags {
			cleanTag := strings.TrimSpace(tag)
			if cleanTag != "" {
				tagsCounts[cleanTag]++
			}
		}
	}

	return models.OverallStats{
		TotalPositions:    len(positions),
		StatusBreakdown:   statusCounts,
		PriorityBreakdown: priorityCounts,
		TagsBreakdown:     tagsCounts,
		RollingBasedCount: rollingCount,
	}
}

func (s *StatsService) GetAnalytics(refreshToken, spreadsheetID, sheetName, groupBy string) (models.AnalyticsResponse, error) {
	positions, err := s.positionService.FetchAllPositions(refreshToken, spreadsheetID, sheetName)
	if err != nil {
		return models.AnalyticsResponse{}, err
	}

	overall := CalculateOverallStats(positions)
	response := models.AnalyticsResponse{Overall: overall}

	groupBy = strings.ToLower(strings.TrimSpace(groupBy))
	if groupBy == "" {
		return response, nil
	}

	if groupBy == "tag" || groupBy == "tags" {
		response.Grouped = groupPositionsByTag(positions)
		return response, nil
	}

	groupedMap := make(map[string][]models.Position)
	for _, pos := range positions {
		key := getGroupKey(pos, groupBy)
		groupedMap[key] = append(groupedMap[key], pos)
	}

	groupedStats := make([]models.GroupedStats, 0, len(groupedMap))
	for key, groupPositions := range groupedMap {
		groupedStats = append(groupedStats, models.GroupedStats{
			GroupKey: key,
			Stats:    CalculateOverallStats(groupPositions),
		})
	}

	response.Grouped = groupedStats
	return response, nil
}

func getGroupKey(pos models.Position, groupBy string) string {
	switch groupBy {
	case "country":
		if pos.University.Country != "" {
			return pos.University.Country
		}
		return "Unknown Country"
	case "university":
		if pos.University.Name != "" {
			return pos.University.Name
		}
		return "Unknown University"
	case "status":
		if pos.Status != "" {
			return string(pos.Status)
		}
		return "Unspecified Status"
	case "priority":
		if pos.Priority != "" {
			return string(pos.Priority)
		}
		return "Unspecified Priority"
	default:
		return "Other"
	}
}

func groupPositionsByTag(positions []models.Position) []models.GroupedStats {
	tagMap := make(map[string][]models.Position)
	for _, pos := range positions {
		for _, tag := range pos.Tags {
			cleanTag := strings.TrimSpace(tag)
			if cleanTag != "" {
				tagMap[cleanTag] = append(tagMap[cleanTag], pos)
			}
		}
	}

	grouped := make([]models.GroupedStats, 0, len(tagMap))
	for tag, posList := range tagMap {
		grouped = append(grouped, models.GroupedStats{
			GroupKey: tag,
			Stats:    CalculateOverallStats(posList),
		})
	}
	return grouped
}
