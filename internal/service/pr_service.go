package service

import (
	"math/rand"
	"time"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type prService struct {
	prs   domain.PullRequestRepository
	users domain.UserRepository
	teams domain.TeamRepository
}

func NewPRService(p domain.PullRequestRepository, u domain.UserRepository, t domain.TeamRepository) domain.PRService {
	return &prService{prs: p, users: u, teams: t}
}

func chooseRandom(arr []string, count int) []string {
	if len(arr) <= count {
		return append([]string{}, arr...)
	}
	perm := rand.Perm(len(arr))
	res := make([]string, 0, count)
	for i := 0; i < count; i++ {
		res = append(res, arr[perm[i]])
	}
	return res
}

// CREATE PR
func (s *prService) CreatePR(id, name, authorID string) (*domain.PullRequest, error) {
	_, err := s.prs.GetPR(id)
	if err == nil {
		return nil, domain.NewError(domain.ErrorPRExists, "PR id already exists")
	}

	author, err := s.users.GetUser(authorID)
	if err != nil {
		return nil, err
	}

	team, err := s.teams.GetTeam(author.TeamName)
	if err != nil {
		return nil, err
	}

	candidates := []string{}
	for _, m := range team.Members {
		if m.IsActive && m.UserID != authorID {
			candidates = append(candidates, m.UserID)
		}
	}

	assigned := chooseRandom(candidates, 2)
	now := time.Now().UTC()

	pr := &domain.PullRequest{
		PullRequestID:   id,
		PullRequestName: name,
		AuthorID:        authorID,
		Status:          domain.StatusOpen,
		Assigned:        assigned,
		CreatedAt:       &now,
	}

	return pr, s.prs.CreatePR(pr)
}

// MERGE PR
func (s *prService) MergePR(id string) (*domain.PullRequest, error) {
	pr, err := s.prs.GetPR(id)
	if err != nil {
		return nil, err
	}

	if pr.Status != domain.StatusMerged {
		now := time.Now().UTC()
		pr.Status = domain.StatusMerged
		pr.MergedAt = &now
		s.prs.UpdatePR(pr)
	}
	return pr, nil
}

// REASSIGN REVIEWER
func (s *prService) ReassignReviewer(prID, oldUserID string) (*domain.PullRequest, string, error) {
	pr, err := s.prs.GetPR(prID)
	if err != nil {
		return nil, "", err
	}

	if pr.Status == domain.StatusMerged {
		return nil, "", domain.NewError(domain.ErrorPRMerged, "PR merged")
	}

	idx := -1
	for i, r := range pr.Assigned {
		if r == oldUserID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, "", domain.NewError(domain.ErrorNotAssigned, "old user not assigned")
	}

	oldUser, err := s.users.GetUser(oldUserID)
	if err != nil {
		return nil, "", err
	}

	team, err := s.teams.GetTeam(oldUser.TeamName)
	if err != nil {
		return nil, "", err
	}

	var candidates []string
	for _, m := range team.Members {
		if !m.IsActive || m.UserID == oldUserID || m.UserID == pr.AuthorID {
			continue
		}
		already := false
		for _, a := range pr.Assigned {
			if a == m.UserID {
				already = true
			}
		}
		if !already {
			candidates = append(candidates, m.UserID)
		}
	}

	if len(candidates) == 0 {
		return nil, "", domain.NewError(domain.ErrorNoCandidate, "no candidate")
	}

	newR := chooseRandom(candidates, 1)[0]
	pr.Assigned[idx] = newR
	s.prs.UpdatePR(pr)

	return pr, newR, nil
}
