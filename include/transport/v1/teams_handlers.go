package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

type TeamService interface{
    CreateTeam(ctx context.Context, req models.Team) (*transport.TeamCreateResponse, error)
    GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

type TeamHandler struct {
	teamService TeamService
}

func NewTeamHandler(teamService TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
    var team models.Team

    if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
        sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
        return
    }

    result_team, err := h.teamService.CreateTeam(r.Context(),team)
    if err != nil {
        handleServiceError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(result_team); err != nil {
        sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
        return
    }
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
    teamName := r.URL.Query().Get("team_name")
    if teamName == "" {
        sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
        return
    }

    team, err := h.teamService.GetTeam(r.Context(), teamName)
    if err != nil {
        handleServiceError(w, err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(team); err != nil {
        sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
        return
    }
}