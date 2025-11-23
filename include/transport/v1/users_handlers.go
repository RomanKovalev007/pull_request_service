package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

type UserService interface{
    GetUserPullRequests(ctx context.Context, userID string) (*transport.UserPRsResponse, error)
    SetUserIsActive(ctx context.Context, req transport.UserSetActiveRequest) (*transport.UserSetActiveResponse, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SetUserIsActive(w http.ResponseWriter, r *http.Request){
    var user_request transport.UserSetActiveRequest

    if err := json.NewDecoder(r.Body).Decode(&user_request); err != nil {
        sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
        return
    }

    result_user, err := h.userService.SetUserIsActive(r.Context(), user_request)
    if err != nil {
        handleServiceError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result_user); err != nil {
        sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
        return
    }
}

func (h *UserHandler) GetUserPullRequests(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
        return
    }

    prs, err := h.userService.GetUserPullRequests(r.Context(), userID)
    if err != nil {
        handleServiceError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(prs); err != nil {
        sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
        return
    }
}