package domain

type PullRequestRepository interface {
	CreatePR(pr *PullRequest) error
	GetPR(id string) (*PullRequest, error)
	UpdatePR(pr *PullRequest) error
	GetPRsAssignedTo(userID string) ([]*PullRequest, error)

	// статистика
	GetAssignmentCountByUser() ([]UserAssignmentStat, error)
	GetReviewerCountByPR() ([]PRAssignmentStat, error)
}
