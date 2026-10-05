package models

type OverallStats struct {
	TotalPositions    int              `json:"total_positions"`
	StatusBreakdown   map[Status]int   `json:"status_breakdown"`   // e.g. {"Accepted": 2, "Rejected": 5, "In progress": 3}
	PriorityBreakdown map[Priority]int `json:"priority_breakdown"` // e.g. {"High": 6, "Medium": 4}
	TagsBreakdown     map[string]int   `json:"tags_breakdown"`     // e.g. {"Machine Learning": 3, "NLP": 2}
	RollingBasedCount int              `json:"rolling_based_count"`
}

type GroupedStats struct {
	GroupKey string       `json:"group_key"`
	Stats    OverallStats `json:"stats"`
}

type AnalyticsResponse struct {
	Overall OverallStats   `json:"overall"`
	Grouped []GroupedStats `json:"grouped,omitempty"`
}
