package service

import (
	"context"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

type Repository interface{
	CreateTeam(ctx context.Context, team models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)

	SetUserIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserPullRequests(ctx context.Context,userID string) ([]models.PullRequestShort, error)

	CreatePullRequest(ctx context.Context, req models.PullRequestShort) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequest, string, error)
}

type Service struct {
    db Repository
}

func NewService(db Repository) *Service {
    return &Service{db: db}
}