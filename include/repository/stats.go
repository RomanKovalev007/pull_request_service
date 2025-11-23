package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

type StatsRepository struct {
    db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
    return &StatsRepository{db: db}
}

func (r *StatsRepository) GetUserStats(ctx context.Context) ([]models.UserStat, error) {
    query := `
        SELECT u.id, u.username, u.team_name, COUNT(pr.reviewer_id) as assignment_count
        FROM users u
        LEFT JOIN pr_reviewers pr ON u.id = pr.reviewer_id
        WHERE u.is_active = true
        GROUP BY u.id, u.username, u.team_name
        ORDER BY assignment_count DESC
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query user stats: %w", err)
    }
    defer rows.Close()

    var stats []models.UserStat
    for rows.Next() {
        var stat models.UserStat
        if err := rows.Scan(&stat.UserID, &stat.Username, &stat.TeamName, &stat.AssignmentCount); err != nil {
            return nil, fmt.Errorf("failed to scan user stat: %w", err)
        }
        stats = append(stats, stat)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return stats, nil
}


func (r *StatsRepository) GetPRStats(ctx context.Context) ([]models.PullRequestStat, error) {
    query := `
        SELECT p.id, p.pull_request_name, p.author_id, p.status, 
			COUNT(pr.reviewer_id) as reviewer_count, p.created_at
        FROM pull_requests p
        LEFT JOIN pr_reviewers pr ON p.id = pr.pull_request_id
        GROUP BY p.id, p.pull_request_name, p.author_id, p.status, p.created_at
        ORDER BY p.created_at DESC
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query PR stats: %w", err)
    }
    defer rows.Close()

    var stats []models.PullRequestStat
    for rows.Next() {
        var stat models.PullRequestStat
        if err := rows.Scan(
            &stat.PullRequestID,
            &stat.PullRequestName,
            &stat.AuthorID,
            &stat.Status,
            &stat.AssignedCount,
            &stat.CreatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan PR stat: %w", err)
        }
        stats = append(stats, stat)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return stats, nil
}


func (r *StatsRepository) GetTeamStats(ctx context.Context) ([]models.TeamStat, error) {
    query := `
		SELECT 
			t.team_name,
			COUNT(DISTINCT u.id) as member_count,
			COUNT(DISTINCT active_reviewers.reviewer_id) as active_reviewers_count,
			COUNT(DISTINCT p.id) as open_prs_count
		FROM teams t
		LEFT JOIN users u ON t.team_name = u.team_name AND u.is_active = true
		LEFT JOIN pull_requests p ON u.id = p.author_id AND p.status = 'OPEN'
		LEFT JOIN (
			SELECT DISTINCT pr.reviewer_id 
			FROM pr_reviewers pr
			JOIN pull_requests p ON pr.pull_request_id = p.id
			WHERE p.status = 'OPEN'
		) active_reviewers ON active_reviewers.reviewer_id = u.id
		GROUP BY t.team_name
		ORDER BY active_reviewers_count DESC
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query team stats: %w", err)
    }
    defer rows.Close()

    var stats []models.TeamStat
    for rows.Next() {
        var stat models.TeamStat
        if err := rows.Scan(
            &stat.TeamName,
            &stat.MemberCount,
            &stat.ActiveReviewers,
            &stat.ActivePRs,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan team stat: %w", err)
        }
        stats = append(stats, stat)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return stats, nil
}