package postgres

import (
	"context"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoPostgres struct {
	db *pgxpool.Pool
}

func NewUserRepoPostgres(db *pgxpool.Pool) *UserRepoPostgres {
	return &UserRepoPostgres{db: db}
}

func (r *UserRepoPostgres) UpsertUser(u *domain.User) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users(user_id, username, team_name, is_active)
         VALUES ($1,$2,$3,$4)
         ON CONFLICT (user_id)
         DO UPDATE SET username=$2, team_name=$3, is_active=$4`,
		u.UserID, u.Username, u.TeamName, u.IsActive,
	)
	return err
}

func (r *UserRepoPostgres) GetUser(id string) (*domain.User, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT user_id, username, team_name, is_active
         FROM users WHERE user_id=$1`,
		id,
	)

	var u domain.User
	err := row.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)
	if err == pgx.ErrNoRows {
		return nil, domain.NewError(domain.ErrorNotFound, "user not found")
	}
	return &u, err
}

func (r *UserRepoPostgres) UpdateUser(u *domain.User) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET username=$2, team_name=$3, is_active=$4 WHERE user_id=$1`,
		u.UserID, u.Username, u.TeamName, u.IsActive,
	)
	return err
}

func (r *UserRepoPostgres) GetUsersByTeam(teamName string) ([]*domain.User, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT user_id, username, team_name, is_active
         FROM users WHERE team_name=$1`,
		teamName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*domain.User
	for rows.Next() {
		var u domain.User
		rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)
		res = append(res, &u)
	}
	return res, nil
}
