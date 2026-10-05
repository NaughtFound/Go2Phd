package models

import "time"

type Priority string
type Status string

const (
	PriorityHigh   Priority = "High"
	PriorityMedium Priority = "Medium"
	PriorityLow    Priority = "Low"
)

const (
	StatusNotStarted  Status = "Not started"
	StatusInProgress  Status = "In progress"
	StatusApplied     Status = "Applied"
	StatusEmailed     Status = "Emailed"
	StatusInterview   Status = "Interview"
	StatusRejected    Status = "Rejected"
	StatusShortlisted Status = "Shortlisted"
	StatusAccepted    Status = "Accepted"
	StatusCanceled    Status = "Canceled"
	StatusNoAnswer    Status = "No Answer"
)

type Position struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Priority     Priority   `json:"priority"`
	Tags         []string   `json:"tags"`
	Status       Status     `json:"status"`
	University   University `json:"university"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	RollingBased bool       `json:"rolling_based"`
	Links        []string   `json:"links"`
}
