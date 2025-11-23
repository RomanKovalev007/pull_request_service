package service

import (
	"context"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

type prRepository interface{
	CreatePullRequest(ctx context.Context, req models.PullRequestShort) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequest, string, error)
}

type PrService struct{
	prRepo prRepository
}

func NewPrService(prRepo prRepository) *PrService{
	return &PrService{prRepo: prRepo}
}

func (s *PrService) CreatePullRequest(ctx context.Context, req transport.CreatePRRequest) (*transport.CreatePRResponse, error) {
	if err := s.validateCreatePR(req); err != nil {
		return nil, err
	}

	req_pr := models.PullRequestShort{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	pr, err := s.prRepo.CreatePullRequest(ctx, req_pr)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to create pull request"}
	}

	resp := transport.CreatePRResponse{PullRequest: *pr}

    return &resp, nil
}

func (s *PrService) MergePullRequest(ctx context.Context, req transport.MergePRRequest) (*transport.MergePRResponse, error) {
	if err := s.validateMergePR(req); err != nil {
		return nil, err
	}

	pr, err := s.prRepo.MergePullRequest(ctx, req.PullRequestID)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to merge pull request"}
	}

	resp := transport.MergePRResponse{PullRequest: *pr}

    return &resp, nil
}

func (s *PrService) ReassignReviewer(ctx context.Context, req transport.ReassignRequest) (*transport.ReassignResponse, error) {
	if err := s.validateReassignReviewer(req); err != nil {
		return nil, err
	}

	pr, newID, err := s.prRepo.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		return nil, &ServiceError{Code: err.Error(), Message: "failed to reassign pull request"}
	}
	
	resp := transport.ReassignResponse{PullRequest: *pr, ReplacedBy: newID}

    return &resp, nil
}