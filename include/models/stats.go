package models

import "time"

type UserStat struct {
	UserID          string `json:"user_id"`
	Username        string `json:"username"`
	TeamName        string `json:"team_name"`
	AssignmentCount int    `json:"assignment_count"`
}

type PullRequestStat struct {
	PullRequestID   string    `json:"pull_request_id"`
	PullRequestName string    `json:"pull_request_name"`
	AuthorID        string    `json:"author_id"`
	Status          string    `json:"status"`
	AssignedCount   int       `json:"assigned_reviewers"`
	CreatedAt       time.Time `json:"created_at"`
}

type TeamStat struct {
	TeamName        string `json:"team_name"`
	MemberCount     int    `json:"member_count"`
	ActiveReviewers int    `json:"active_reviewers"`
	ActivePRs       int    `json:"active_prs"`
}

type TotalStats struct {
	TotalUsers           int `json:"total_users"`
	TotalPRs             int `json:"total_prs"`
	TotalTeams           int `json:"total_teams"`
	TotalActiveReviewers int `json:"total_active_reviewers"`
	MergedPRs            int `json:"merged_prs"`
	OpenPRs              int `json:"open_prs"`
}

type StatsResponse struct {
	UserStats  []UserStat        `json:"user_stats"`
	PRStats    []PullRequestStat `json:"pr_stats"`
	TeamStats  []TeamStat        `json:"team_stats"`
	TotalStats TotalStats        `json:"total_stats"`
	Timestamp  time.Time         `json:"timestamp"`
}
