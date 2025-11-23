package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

type UserRepository struct{
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error){
	var user models.User 

	err := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id=$2
		RETURNING id, username, team_name, is_active`,
        isActive, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if err != nil{
		if err == sql.ErrNoRows{
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to set user isActive status: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserPullRequests(ctx context.Context,userID string) ([]models.PullRequestShort, error) {
	tx, err := r.db.Begin()
    if err != nil {
        return nil, fmt.Errorf("failed to begin tx: %w", err)
    }

    var exists bool
    err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
    if err != nil {
        return nil, fmt.Errorf("failed to check user exists: %w", err)
    }
    if !exists {
        return nil, ErrNotFound
    }

    var prs []models.PullRequestShort

    rows, err := tx.QueryContext(ctx,`
        SELECT p.id, p.pull_request_name, p.author_id, p.status
        FROM pull_requests p
        JOIN pr_reviewers pr ON p.id = pr.pull_request_id
        WHERE pr.reviewer_id = $1
        ORDER BY p.created_at DESC`,
        userID)
    if err != nil {
        return nil, fmt.Errorf("failed to select user pull requests: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var pr models.PullRequestShort
        if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
            return nil, fmt.Errorf("failed to scan user pull requests: %w", err)
        }
        prs = append(prs, pr)
    }

	err = tx.Commit()
    if err != nil {
        return nil, fmt.Errorf("failed to tx commit: %w", err)
    }

    return prs, nil
}