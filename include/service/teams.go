package service

import (
	"context"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

type teamRepository interface {
	CreateTeam(ctx context.Context, team models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

type TeamService struct {
	teamRepo teamRepository
}

func NewTeamService(teamRepo teamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, req models.Team) (*transport.TeamCreateResponse, error) {
	if err := s.validateCreateTeam(req); err != nil {
		return nil, err
	}

	team, err := s.teamRepo.CreateTeam(ctx, req)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to create team"}
	}

	return &transport.TeamCreateResponse{Team: *team}, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	if err := s.validateGetTeam(teamName); err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to get team"}
	}

	team, err := s.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to get team"}
	}
	return team, nil
}
