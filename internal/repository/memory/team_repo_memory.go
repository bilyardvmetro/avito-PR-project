package memory

import (
	"errors"
	"sync"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type TeamRepoMemory struct {
	mu    sync.RWMutex
	teams map[string]*domain.Team
}

func NewTeamRepoMemory() *TeamRepoMemory {
	return &TeamRepoMemory{
		teams: make(map[string]*domain.Team),
	}
}

func (r *TeamRepoMemory) CreateTeam(t *domain.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.teams[t.TeamName]; ok {
		return errors.New("exists")
	}
	r.teams[t.TeamName] = t
	return nil
}

func (r *TeamRepoMemory) GetTeam(name string) (*domain.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.teams[name]
	if !ok {
		return nil, domain.NewError(domain.ErrorNotFound, "team not found")
	}
	return t, nil
}

func (r *TeamRepoMemory) UpdateTeam(t *domain.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.teams[t.TeamName] = t
	return nil
}
