package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	"github.com/RomanKovalev007/pull_request_service/include/service"
)

type StatsService interface {
	GetStats(ctx context.Context) (*models.StatsResponse, error)
}

type StatsHandler struct {
	statsService StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.statsService.GetStats(r.Context())
	if err != nil {
		http.Error(w, "Failed to get statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
