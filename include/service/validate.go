package service

import (
	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

func (s *TeamService) validateCreateTeam(team models.Team) *ServiceError {
	if team.TeamName == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "team_name is required"}
	}

	if len(team.Members) == 0 {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "team must have at least one member"}
	}

	userIDs := make(map[string]bool)
	for _, member := range team.Members {
		if member.UserID == "" {
			return &ServiceError{Code: ErrInvalidInput.Error(), Message: "user_id is required for all members"}
		}
		if member.Username == "" {
			return &ServiceError{Code: ErrInvalidInput.Error(), Message: "username is required for all members"}
		}
		if userIDs[member.UserID] {
			return &ServiceError{Code: ErrInvalidInput.Error(), Message: "duplicate user_id in team members"}
		}
		userIDs[member.UserID] = true
	}

	return nil
}

func (s *TeamService) validateGetTeam(teamName string) *ServiceError {
	if teamName == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "team_name is required"}
	}
	return nil
}

func (s *UserService) validateSetUserIsActive(req transport.UserSetActiveRequest) *ServiceError {
	if req.UserID == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "user_id is required"}
	}
	return nil
}

func (s *UserService) validateGetUserPullRequests(userID string) *ServiceError {
	if userID == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "user_id is required"}
	}
	return nil
}

func (s *PrService) validateCreatePR(req transport.CreatePRRequest) *ServiceError {
	if req.PullRequestID == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "pull_request_id is required"}
	}
	if req.PullRequestName == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "pull_request_name is required"}
	}
	if req.AuthorID == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "author_id is required"}
	}

	return nil
}

func (s *PrService) validateMergePR(req transport.MergePRRequest) *ServiceError {
	if req.PullRequestID == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "pull_request_id is required"}
	}

	return nil
}

func (s *PrService) validateReassignReviewer(req transport.ReassignRequest) *ServiceError {
	if req.PullRequestID == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "pull_request_id is required"}
	}

	if req.OldUserID == "" {
		return &ServiceError{Code: ErrInvalidInput.Error(), Message: "old_user_id is required"}
	}

	return nil
}
