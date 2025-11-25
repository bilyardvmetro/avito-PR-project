package service

import (
	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type teamService struct {
	teams domain.TeamRepository
	users domain.UserRepository
}

func NewTeamService(t domain.TeamRepository, u domain.UserRepository) domain.TeamService {
	return &teamService{teams: t, users: u}
}

func (s *teamService) AddTeam(team domain.Team) (*domain.Team, error) {
	if _, err := s.teams.GetTeam(team.TeamName); err == nil {
		return nil, domain.NewError(domain.ErrorTeamExists, "team already exists")
	}

	for _, m := range team.Members {
		s.users.UpsertUser(&domain.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: team.TeamName,
			IsActive: m.IsActive,
		})
	}

	err := s.teams.CreateTeam(&team)
	return &team, err
}

func (s *teamService) GetTeam(name string) (*domain.Team, error) {
	return s.teams.GetTeam(name)
}
