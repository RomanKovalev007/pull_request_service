package v1

import (
	"encoding/json"
	"net/http"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

type TeamHandler struct{
	service Service
}

func NewTeamHandler(service Service) *TeamHandler{
	return &TeamHandler{service: service}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
    var team models.Team

    if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
        sendError(w, http.StatusBadRequest, models.INVALID_INPUT, "invalid request payload")
        return
    }

    result_team, err := h.service.CreateTeam(r.Context(),team)
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

    team, err := h.service.GetTeam(r.Context(), teamName)
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