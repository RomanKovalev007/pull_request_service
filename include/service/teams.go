package service

import (
	"context"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)


func (s *Service) CreateTeam(ctx context.Context, req models.Team) (*transport.TeamCreateResponse, error) {
    if err := s.validateCreateTeam(req); err != nil{
		return nil, err
	}

	team, err := s.db.CreateTeam(ctx, req)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to create team"}
	}

    return &transport.TeamCreateResponse{Team: *team}, nil
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	s.validateGetTeam(teamName)

	team, err := s.db.GetTeam(ctx, teamName)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to get team"}
	}
    return team, nil
}