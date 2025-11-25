package domain

type TeamRepository interface {
	CreateTeam(team *Team) error
	GetTeam(name string) (*Team, error)
	UpdateTeam(team *Team) error
}
