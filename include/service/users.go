package service

import (
	"context"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

type userRepository interface{
	SetUserIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserPullRequests(ctx context.Context,userID string) ([]models.PullRequestShort, error)
}

type UserService struct{
	userRepo userRepository
}

func NewUserService(userRepo userRepository) *UserService{
	return &UserService{userRepo: userRepo}
}

func (s *UserService) SetUserIsActive(ctx context.Context, req transport.UserSetActiveRequest) (*transport.UserSetActiveResponse, error) {
	s.validateSetUserIsActive(req)

	user, err := s.userRepo.SetUserIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to set user active status"}
	}

	resp := transport.UserSetActiveResponse{
		User: *user,
	}

    return &resp, nil
}

func (s *UserService) GetUserPullRequests(ctx context.Context, userID string) (*transport.UserPRsResponse, error) {
    s.validateGetUserPullRequests(userID)

	prs, err := s.userRepo.GetUserPullRequests(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to get user pull requests"}
	}

	resp := transport.UserPRsResponse{
		UserID: userID,
		PullRequests: prs,
	}

    return &resp, nil
}