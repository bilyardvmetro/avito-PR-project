package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
)

type TeamRepoPostgres struct {
	db *pgxpool.Pool
}

func NewTeamRepoPostgres(db *pgxpool.Pool) *TeamRepoPostgres {
	return &TeamRepoPostgres{db: db}
}

func (r *TeamRepoPostgres) CreateTeam(team *domain.Team) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO teams(team_name) VALUES ($1)`,
		team.TeamName,
	)
	return err
}

func (r *TeamRepoPostgres) GetTeam(name string) (*domain.Team, error) {
	row := r.db.QueryRow(
		context.Background(),
		`SELECT team_name FROM teams WHERE team_name = $1`,
		name,
	)

	var t domain.Team
	err := row.Scan(&t.TeamName)
	if err == pgx.ErrNoRows {
		return nil, domain.NewError(domain.ErrorNotFound, "team not found")
	}
	if err != nil {
		return nil, err
	}

	// Members здесь можем оставить пустым — TeamService при желании
	// может дополнительно дотянуть юзеров из UserRepository.
	t.Members = nil
	return &t, nil
}

func (r *TeamRepoPostgres) UpdateTeam(team *domain.Team) error {
	// пока ничего не делаем, по ТЗ апдейта команды нет
	return nil
}
