package service

import (
	"context"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

func (s *Service) CreatePullRequest(ctx context.Context, req transport.CreatePRRequest) (*transport.CreatePRResponse, error) {
	if err := s.validateCreatePR(req); err != nil {
		return nil, err
	}

	req_pr := models.PullRequestShort{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	pr, err := s.db.CreatePullRequest(ctx, req_pr)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to create pull request"}
	}

	resp := transport.CreatePRResponse{PullRequest: *pr}

    return &resp, nil
}

func (s *Service) MergePullRequest(ctx context.Context, req transport.MergePRRequest) (*transport.MergePRResponse, error) {
	if err := s.validateMergePR(req); err != nil {
		return nil, err
	}

	pr, err := s.db.MergePullRequest(ctx, req.PullRequestID)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to merge pull request"}
	}

	resp := transport.MergePRResponse{PullRequest: *pr}

    return &resp, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, req transport.ReassignRequest) (*transport.ReassignResponse, error) {
	if err := s.validateReassignReviewer(req); err != nil {
		return nil, err
	}

	pr, newID, err := s.db.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to reassign pull request"}
	}
	
	resp := transport.ReassignResponse{PullRequest: *pr, ReplacedBy: newID}

    return &resp, nil
}