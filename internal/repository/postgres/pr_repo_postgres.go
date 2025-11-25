package postgres

import (
	"context"

	"github.com/bilyardvmetro/avito-PR-project/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PRRepoPostgres struct {
	db *pgxpool.Pool
}

func NewPRRepoPostgres(db *pgxpool.Pool) *PRRepoPostgres {
	return &PRRepoPostgres{db: db}
}

func (r *PRRepoPostgres) CreatePR(pr *domain.PullRequest) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO pull_requests
         (pull_request_id, pull_request_name, author_id, status, created_at)
         VALUES ($1,$2,$3,$4,$5)`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.CreatedAt,
	)
	if err != nil {
		return err
	}

	for _, rev := range pr.Assigned {
		_, err := r.db.Exec(context.Background(),
			`INSERT INTO pr_reviewers(pull_request_id, reviewer_id)
			 VALUES ($1,$2)`,
			pr.PullRequestID, rev,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PRRepoPostgres) GetPR(id string) (*domain.PullRequest, error) {
	var pr domain.PullRequest

	row := r.db.QueryRow(context.Background(),
		`SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
         FROM pull_requests WHERE pull_request_id=$1`,
		id,
	)

	err := row.Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.NewError(domain.ErrorNotFound, "PR not found")
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		`SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id=$1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rID string
		rows.Scan(&rID)
		pr.Assigned = append(pr.Assigned, rID)
	}

	return &pr, nil
}

func (r *PRRepoPostgres) UpdatePR(pr *domain.PullRequest) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE pull_requests
         SET pull_request_name=$2, author_id=$3, status=$4, created_at=$5, merged_at=$6
         WHERE pull_request_id=$1`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt,
	)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(context.Background(),
		`DELETE FROM pr_reviewers WHERE pull_request_id=$1`,
		pr.PullRequestID,
	)

	for _, rev := range pr.Assigned {
		_, err = r.db.Exec(context.Background(),
			`INSERT INTO pr_reviewers(pull_request_id, reviewer_id)
			 VALUES ($1,$2)`,
			pr.PullRequestID, rev,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PRRepoPostgres) GetPRsAssignedTo(userID string) ([]*domain.PullRequest, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT p.pull_request_id, p.pull_request_name, p.author_id, p.status
         FROM pull_requests p
         JOIN pr_reviewers r
         ON p.pull_request_id = r.pull_request_id
         WHERE r.reviewer_id=$1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
		res = append(res, &pr)
	}

	return res, nil
}

// статистика
func (r *PRRepoPostgres) GetAssignmentCountByUser() ([]domain.UserAssignmentStat, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT reviewer_id, COUNT(*) 
         FROM pr_reviewers
         GROUP BY reviewer_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.UserAssignmentStat
	for rows.Next() {
		var s domain.UserAssignmentStat
		if err := rows.Scan(&s.UserID, &s.Assignments); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (r *PRRepoPostgres) GetReviewerCountByPR() ([]domain.PRAssignmentStat, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT pull_request_id, COUNT(*) 
         FROM pr_reviewers
         GROUP BY pull_request_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.PRAssignmentStat
	for rows.Next() {
		var s domain.PRAssignmentStat
		if err := rows.Scan(&s.PullRequestID, &s.Reviewers); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}
