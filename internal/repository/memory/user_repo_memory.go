package memory

import (
	"sync"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type UserRepoMemory struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepoMemory() *UserRepoMemory {
	return &UserRepoMemory{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepoMemory) UpsertUser(u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *u
	r.users[u.UserID] = &copy
	return nil
}

func (r *UserRepoMemory) GetUser(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, domain.NewError(domain.ErrorNotFound, "user not found")
	}
	copy := *u
	return &copy, nil
}

func (r *UserRepoMemory) UpdateUser(u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[u.UserID]; !ok {
		return domain.NewError(domain.ErrorNotFound, "user not found")
	}
	copy := *u
	r.users[u.UserID] = &copy
	return nil
}

func (r *UserRepoMemory) GetUsersByTeam(teamName string) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var res []*domain.User
	for _, u := range r.users {
		if u.TeamName == teamName {
			copy := *u
			res = append(res, &copy)
		}
	}
	return res, nil
}
