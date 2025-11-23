package v1

import (
	"encoding/json"
	"net/http"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	"github.com/RomanKovalev007/pull_request_service/include/repository"
	"github.com/RomanKovalev007/pull_request_service/include/service"
)

func sendError(w http.ResponseWriter, statusCode int, errorCode models.ErrorResponseErrorCode, message string) {
	var errResp models.ErrorResponse
	errResp.Error.Code = errorCode
	errResp.Error.Message = message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		return
	}
}

func handleServiceError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(*service.ServiceError); ok {
		switch serviceErr.Code {
		case repository.ErrTeamExists.Error():
			sendError(w, http.StatusBadRequest, models.TEAMEXISTS, serviceErr.Message)
		case repository.ErrPRExists.Error():
			sendError(w, http.StatusConflict, models.PREXISTS, serviceErr.Message)
		case repository.ErrPRMerged.Error():
			sendError(w, http.StatusConflict, models.PRMERGED, serviceErr.Message)
		case repository.ErrNotAssigned.Error():
			sendError(w, http.StatusConflict, models.NOTASSIGNED, serviceErr.Message)
		case repository.ErrNoCandidate.Error():
			sendError(w, http.StatusConflict, models.NOCANDIDATE, serviceErr.Message)
		case repository.ErrNotFound.Error():
			sendError(w, http.StatusNotFound, models.NOTFOUND, serviceErr.Message)
		case service.ErrInvalidInput.Error():
			sendError(w, http.StatusBadRequest, models.INVALID_INPUT, serviceErr.Message)
		default:
			sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, serviceErr.Message)
		}
	} else {
		sendError(w, http.StatusInternalServerError, models.INTERNAL_ERROR, "internal server error")
	}
}
