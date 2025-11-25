package service

import "github.com/bilyardvmetro/avito-PR-project/internal/domain"

type statsService struct {
	prs domain.PullRequestRepository
}

func NewStatsService(prs domain.PullRequestRepository) domain.StatsService {
	return &statsService{prs: prs}
}

func (s *statsService) GetAssignmentStats() (*domain.AssignmentStats, error) {
	byUser, err := s.prs.GetAssignmentCountByUser()
	if err != nil {
		return nil, err
	}

	byPR, err := s.prs.GetReviewerCountByPR()
	if err != nil {
		return nil, err
	}

	return &domain.AssignmentStats{
		ByUser: byUser,
		ByPR:   byPR,
	}, nil
}
