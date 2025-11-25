package domain

import "time"

const (
	StatusOpen   = "OPEN"
	StatusMerged = "MERGED"
)

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequest struct {
	PullRequestID   string     `json:"pull_request_id"`
	PullRequestName string     `json:"pull_request_name"`
	AuthorID        string     `json:"author_id"`
	Status          string     `json:"status"`
	Assigned        []string   `json:"assigned_reviewers"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	MergedAt        *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// статистика
type UserAssignmentStat struct {
	UserID      string `json:"user_id"`
	Assignments int64  `json:"assignments"`
}

type PRAssignmentStat struct {
	PullRequestID string `json:"pull_request_id"`
	Reviewers     int64  `json:"reviewers"`
}

type AssignmentStats struct {
	ByUser []UserAssignmentStat `json:"by_user"`
	ByPR   []PRAssignmentStat   `json:"by_pr"`
}
