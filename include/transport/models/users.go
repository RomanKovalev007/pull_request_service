package transport

import "github.com/RomanKovalev007/pull_request_service/include/models"

type UserSetActiveRequest struct {
    UserID   string `json:"user_id"`
    IsActive bool   `json:"is_active"`
}

type UserSetActiveResponse struct {
    User models.User `json:"user"`
}

type UserPRsResponse struct {
    UserID        string               `json:"user_id"`
    PullRequests  []models.PullRequestShort   `json:"pull_requests"`
}