package service

import (
	"context"

	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

func (s *Service) SetUserIsActive(ctx context.Context, req transport.UserSetActiveRequest) (*transport.UserSetActiveResponse, error) {
	s.validateSetUserIsActive(req)

	user, err := s.db.SetUserIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to set user active status"}
	}

	resp := transport.UserSetActiveResponse{
		User: *user,
	}

    return &resp, nil
}

func (s *Service) GetUserPullRequests(ctx context.Context, userID string) (*transport.UserPRsResponse, error) {
    s.validateGetUserPullRequests(userID)

	prs, err := s.db.GetUserPullRequests(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to get user pull requests"}
	}

	resp := transport.UserPRsResponse{
		UserID: userID,
		PullRequests: prs,
	}

    return &resp, nil
}