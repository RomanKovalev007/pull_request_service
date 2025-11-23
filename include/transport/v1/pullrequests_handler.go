package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

type PRService interface {
	CreatePullRequest(ctx context.Context, req transport.CreatePRRequest) (*transport.CreatePRResponse, error)
	MergePullRequest(ctx context.Context, req transport.MergePRRequest) (*transport.MergePRResponse, error)
	ReassignReviewer(ctx context.Context, req transport.ReassignRequest) (*transport.ReassignResponse, error)
}

type PRHandler struct {
	prService PRService
}

func NewPRHandler(prService PRService) *PRHandler {
	return &PRHandler{prService: prService}
}

func (h *PRHandler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req transport.CreatePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
		return
	}

	pr, err := h.prService.CreatePullRequest(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pr); err != nil {
		sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
		return
	}
}

func (h *PRHandler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req transport.MergePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
		return
	}

	pr, err := h.prService.MergePullRequest(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pr); err != nil {
		sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
		return
	}
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req transport.ReassignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
		return
	}

	pr, err := h.prService.ReassignReviewer(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pr); err != nil {
		sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
		return
	}
}
