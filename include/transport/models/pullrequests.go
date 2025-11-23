package transport

import "github.com/RomanKovalev007/pull_request_service/include/models"

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_reviewer_id"`
}

type CreatePRResponse struct {
	PullRequest models.PullRequest `json:"pr"`
}

type MergePRResponse struct {
	PullRequest models.PullRequest `json:"pr"`
}

type ReassignResponse struct {
	PullRequest models.PullRequest `json:"pr"`
	ReplacedBy  string             `json:"replaced_by"`
}
