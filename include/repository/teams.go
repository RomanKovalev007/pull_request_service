package repository

import (
	"context"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

func (r *Repo) CreateTeam(ctx context.Context, team models.Team) (*models.Team, error) {
    tx, err := r.DB.Begin()
    if err != nil {
        return nil, err
    }

    defer tx.Rollback()

    var exists bool
    err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)", team.TeamName).Scan(&exists)
    if err != nil {
        return nil, err
    }

    if exists {
        return nil, ErrTeamExists
    }

	var result_team models.Team

    err = tx.QueryRowContext(ctx, "INSERT INTO teams (team_name) VALUES ($1) RETURNING team_name", team.TeamName).Scan(&result_team.TeamName)
    if err != nil {
        return nil, err
    }

    for _, member := range team.Members {
        var result_member models.TeamMember
        err = tx.QueryRowContext(ctx,`
            INSERT INTO users (id, username, team_name, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) 
            DO UPDATE SET username = $2, team_name = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
            RETURNING id, username, is_active`,
            member.UserID, member.Username, team.TeamName, member.IsActive).Scan(&result_member.UserID, &result_member.Username, &result_member.IsActive)
        if err != nil {
            return nil, err
        }
        result_team.Members = append(result_team.Members, result_member)
    }
    err = tx.Commit()
    if err != nil {
        return nil, err
    }
    return &result_team, nil
}

func (r *Repo) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
    var team models.Team
    team.TeamName = teamName

    rows, err := r.DB.QueryContext(ctx, `
        SELECT id, username, is_active 
        FROM users 
        WHERE team_name = $1 
        ORDER BY id`, teamName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var member models.TeamMember
        if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
            return nil, err
        }
        team.Members = append(team.Members, member)
    }

    if len(team.Members) == 0 {
        return nil, ErrNotFound
    }

    return &team, nil
}