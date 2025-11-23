package repository

import (
	"context"
	"database/sql"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

func (r *Repo) CreatePullRequest(ctx context.Context, req models.PullRequestShort) (*models.PullRequest, error) {
    tx, err := r.DB.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    var exists bool
    err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM pull_requests 
		WHERE id = $1)`,
		req.PullRequestID).Scan(&exists)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrPRExists
    }

    var authorTeam string
    var authorActive bool
    err = tx.QueryRowContext(ctx, `
		SELECT team_name, is_active FROM users 
		WHERE id = $1`, 
        req.AuthorID).Scan(&authorTeam, &authorActive)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    } else if err != nil {
        return nil, err
    }

    if !authorActive {
        return nil, ErrNotFound
    }

    rows, err := tx.Query(`
        SELECT id FROM users 
        WHERE team_name = $1 AND is_active = true AND id != $2 
        ORDER BY RANDOM() 
        LIMIT 2`,
        authorTeam, req.AuthorID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var reviewers []string
    for rows.Next() {
        var reviewerID string
        if err := rows.Scan(&reviewerID); err != nil {
            return nil, err
        }
        reviewers = append(reviewers, reviewerID)
    }

	if len(reviewers) == 0{
		return nil, ErrNoCandidate
	}

	var pr models.PullRequest
	pr.AssignedReviewers = reviewers

    err = tx.QueryRowContext(ctx,`
        INSERT INTO pull_requests (id, pull_request_name, author_id) 
        VALUES ($1, $2, $3)
		RETURNING id, pull_request_name, author_id, status`,
        req.PullRequestID, req.PullRequestName, req.AuthorID).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)
    if err != nil {
        return nil, err
    }


    for _, reviewerID := range reviewers {
        _, err = tx.ExecContext(ctx, `
            INSERT INTO pr_reviewers (pull_request_id, reviewer_id) 
            VALUES ($1, $2)`,
            req.PullRequestID, reviewerID)
        if err != nil {
            return nil, err
        }
    }
    

    if err = tx.Commit(); err != nil {
        return nil, err
    }

    return &pr, nil
}

func (r *Repo) MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error) {
    tx, err := r.DB.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    var pr models.PullRequest
    var mergedAt sql.NullTime
    
    err = tx.QueryRowContext(ctx,`
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = CURRENT_TIMESTAMP 
        WHERE id = $1 
        RETURNING id, pull_request_name, author_id, status, merged_at`,
        prID).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.MergedAt)

    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    } else if err != nil {
        return nil, err
    }

    if mergedAt.Valid {
        pr.MergedAt = &mergedAt.Time
    }

    rows, err := tx.QueryContext(ctx, `
        SELECT reviewer_id FROM pr_reviewers 
        WHERE pull_request_id = $1`, prID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var reviewerID string
        if err := rows.Scan(&reviewerID); err != nil {
            return nil, err
        }
        pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
    }

    if err = tx.Commit(); err != nil {
        return nil, err
    }

    return &pr, nil
}

func (r *Repo) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequest, string, error) {
    tx, err := r.DB.Begin()
    if err != nil {
        return nil, "", err
    }
    defer tx.Rollback()

    var status string
    err = tx.QueryRowContext(ctx, `
		SELECT status FROM pull_requests
		WHERE id = $1`,
		prID).Scan(&status)
    if err == sql.ErrNoRows {
        return nil, "", ErrNotFound
    } else if err != nil {
        return nil, "", err
    }

    if status == "MERGED" {
        return nil, "", ErrPRMerged
    }

    var isAssigned bool
    err = tx.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM pr_reviewers
		WHERE pull_request_id = $1 AND reviewer_id = $2)`,
        prID, oldUserID).Scan(&isAssigned)
    if err != nil {
        return nil, "", err
    }
    if !isAssigned {
        return nil, "", ErrNotAssigned
    }

    var teamName string
    err = tx.QueryRowContext(ctx, `
		SELECT team_name FROM users
		WHERE id = $1 AND is_active = true`,
		oldUserID).Scan(&teamName)
    if err == sql.ErrNoRows {
        return nil, "", ErrNoCandidate
    } else if err != nil {
        return nil, "", err
    }

    var newReviewerID string
    err = tx.QueryRowContext(ctx, `
        SELECT u.id FROM users u
        WHERE u.team_name = $1 
        AND u.is_active = true
        AND u.id != (SELECT author_id FROM pull_requests WHERE id = $2)
        AND u.id NOT IN (SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $2)
        ORDER BY RANDOM() 
        LIMIT 1`,
        teamName, prID).Scan(&newReviewerID)

    if err == sql.ErrNoRows {
        return nil, "", ErrNoCandidate
    } else if err != nil {
        return nil, "", err
    }

    _, err = tx.ExecContext(ctx, `
        UPDATE pr_reviewers 
        SET reviewer_id = $1 
        WHERE pull_request_id = $2 AND reviewer_id = $3`,
        newReviewerID, prID, oldUserID)
    if err != nil {
        return nil, "", err
    }

    var pr models.PullRequest
    err = tx.QueryRowContext(ctx, `
        SELECT id, pull_request_name, author_id, status, created_at 
        FROM pull_requests
		WHERE id = $1`, prID).
        Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt)
    if err != nil {
        return nil, "", err
    }

    rows, err := tx.QueryContext(ctx, `
        SELECT reviewer_id FROM pr_reviewers WHERE pull_request_id = $1`, prID)
    if err != nil {
        return nil, "", err
    }
    defer rows.Close()

    for rows.Next() {
        var reviewerID string
        if err := rows.Scan(&reviewerID); err != nil {
            return nil, "", err
        }
        pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
    }

    if err = tx.Commit(); err != nil {
        return nil, "", err
    }

    return &pr, newReviewerID, nil
}