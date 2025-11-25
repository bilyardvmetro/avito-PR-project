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
	} else if derr, ok := err.(*domain.DomainError); ok && derr.Code != domain.ErrorNotFound {
		return nil, err
	}

	if err := s.teams.CreateTeam(&team); err != nil {
		return nil, err
	}

	members := make([]domain.TeamMember, 0, len(team.Members))
	for _, m := range team.Members {
		u := &domain.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: team.TeamName,
			IsActive: m.IsActive,
		}

		if err := s.users.UpsertUser(u); err != nil {
			return nil, err
		}

		members = append(members, domain.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	team.Members = members
	return &team, nil
}

func (s *teamService) GetTeam(name string) (*domain.Team, error) {
	_, err := s.teams.GetTeam(name)
	if err != nil {
		return nil, err
	}

	users, err := s.users.GetUsersByTeam(name)
	if err != nil {
		return nil, err
	}

	members := make([]domain.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, domain.TeamMember{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	return &domain.Team{
		TeamName: name,
		Members:  members,
	}, nil
}
