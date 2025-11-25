package service

import (
	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type userService struct {
	users domain.UserRepository
	prs   domain.PullRequestRepository
}

func NewUserService(u domain.UserRepository, p domain.PullRequestRepository) domain.UserService {
	return &userService{users: u, prs: p}
}

func (s *userService) SetIsActive(userID string, active bool) (*domain.User, error) {
	user, err := s.users.GetUser(userID)
	if err != nil {
		return nil, err
	}
	user.IsActive = active
	err = s.users.UpdateUser(user)
	return user, err
}

func (s *userService) GetAssignedPRs(userID string) ([]domain.PullRequestShort, error) {
	prs, _ := s.prs.GetPRsAssignedTo(userID)

	out := make([]domain.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		out = append(out, domain.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}
	return out, nil
}
