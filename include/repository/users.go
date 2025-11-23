package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserPullRequests(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
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

	rows, err := tx.QueryContext(ctx, `
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

func (r *UserRepository) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user models.User

	err = tx.QueryRowContext(ctx, `
        UPDATE users
        SET is_active = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
        RETURNING id, username, team_name, is_active`,
		isActive, userID).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to set user isActive status: %w", err)
	}

	if !isActive {
		err := r.reassignUserReviews(ctx, tx, userID, user.TeamName)
		if err != nil {
			return nil, fmt.Errorf("failed to reassign user reviews: %w", err)
		}

	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) reassignUserReviews(ctx context.Context, tx *sql.Tx, userID, teamName string) error {

	prsToReassign, err := r.findUserOpenPRs(ctx, tx, userID)
	if err != nil {
		return err
	}

	if len(prsToReassign) == 0 {
		return nil
	}

	for _, pr := range prsToReassign {
		newReviewer, err := r.findReplacementReviewer(ctx, tx, userID, pr.AuthorID, teamName, pr.PRID)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := r.removeReviewer(ctx, tx, pr.PRID, userID); err != nil {
					return err
				}
				continue
			}
			return err
		}
		if err := r.replaceReviewer(ctx, tx, pr.PRID, userID, newReviewer); err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) findUserOpenPRs(ctx context.Context, tx *sql.Tx, userID string) ([]models.PRReview, error) {
	query := `
        SELECT 
            pr.pull_request_id,
            p.author_id
        FROM pr_reviewers pr
        JOIN pull_requests p ON pr.pull_request_id = p.id
        WHERE pr.reviewer_id = $1 
        AND p.status = 'OPEN'
        ORDER BY pr.pull_request_id
    `

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []models.PRReview
	for rows.Next() {
		var pr models.PRReview
		if err := rows.Scan(&pr.PRID, &pr.AuthorID); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, rows.Err()
}

func (r *UserRepository) findReplacementReviewer(ctx context.Context, tx *sql.Tx, oldReviewerID, authorID, teamName, prID string) (string, error) {
	var newReviewer string

	err := tx.QueryRowContext(ctx, `
        SELECT u.id 
        FROM users u
        WHERE u.team_name = $1 
        AND u.is_active = true 
        AND u.id != $2  -- исключаем старого ревьювера
        AND u.id != $3  -- исключаем автора PR
        AND u.id NOT IN (
            SELECT reviewer_id 
            FROM pr_reviewers 
            WHERE pull_request_id = $4
        )
        ORDER BY RANDOM()
        LIMIT 1
    `, teamName, oldReviewerID, authorID, prID).Scan(&newReviewer)

	if err != nil {
		return "", err
	}

	return newReviewer, nil
}

func (r *UserRepository) replaceReviewer(ctx context.Context, tx *sql.Tx, prID, oldReviewerID, newReviewerID string) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE pr_reviewers 
        SET reviewer_id = $1 
        WHERE pull_request_id = $2 AND reviewer_id = $3
    `, newReviewerID, prID, oldReviewerID)
	return err
}

func (r *UserRepository) removeReviewer(ctx context.Context, tx *sql.Tx, prID, reviewerID string) error {
	_, err := tx.ExecContext(ctx, `
        DELETE FROM pr_reviewers 
        WHERE pull_request_id = $1 AND reviewer_id = $2
    `, prID, reviewerID)
	return err
}
