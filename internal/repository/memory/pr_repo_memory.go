package memory

import (
	"sync"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type PRRepoMemory struct {
	mu  sync.RWMutex
	prs map[string]*domain.PullRequest
}

func NewPRRepoMemory() *PRRepoMemory {
	return &PRRepoMemory{
		prs: make(map[string]*domain.PullRequest),
	}
}

func (r *PRRepoMemory) CreatePR(pr *domain.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.prs[pr.PullRequestID]; ok {
		return domain.NewError(domain.ErrorPRExists, "PR already exists")
	}
	copy := *pr
	r.prs[pr.PullRequestID] = &copy
	return nil
}

func (r *PRRepoMemory) GetPR(id string) (*domain.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pr, ok := r.prs[id]
	if !ok {
		return nil, domain.NewError(domain.ErrorNotFound, "PR not found")
	}
	copy := *pr
	return &copy, nil
}

func (r *PRRepoMemory) UpdatePR(pr *domain.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.prs[pr.PullRequestID]; !ok {
		// можно либо возвращать ошибку, либо создать — но по логике сервиса
		// Update вызывается только для уже существующих PR
		return domain.NewError(domain.ErrorNotFound, "PR not found")
	}
	copy := *pr
	r.prs[pr.PullRequestID] = &copy
	return nil
}

func (r *PRRepoMemory) GetPRsAssignedTo(userID string) ([]*domain.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var res []*domain.PullRequest
	for _, pr := range r.prs {
		for _, reviewer := range pr.Assigned {
			if reviewer == userID {
				copy := *pr
				res = append(res, &copy)
				break
			}
		}
	}
	return res, nil
}
