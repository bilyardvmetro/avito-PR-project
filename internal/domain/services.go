package domain

type TeamService interface {
	AddTeam(team Team) (*Team, error)
	GetTeam(name string) (*Team, error)
}

type UserService interface {
	SetIsActive(userID string, active bool) (*User, error)
	GetAssignedPRs(userID string) ([]PullRequestShort, error)
}

type PRService interface {
	CreatePR(id, name, authorID string) (*PullRequest, error)
	MergePR(id string) (*PullRequest, error)
	ReassignReviewer(prID, oldUserID string) (*PullRequest, string, error)
}
