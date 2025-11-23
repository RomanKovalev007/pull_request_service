package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RomanKovalev007/pull_request_service/include/repository"
)

type Health struct {
	Status            healthStatus
	Err               string
	Timestamp         string
	Database          databaseStatus
	Migration_version uint
	Migration_dirty   bool
}

const (
	STATUSHEALTHY     healthStatus   = "healthy"
	STATUSUNHEALTHY   healthStatus   = "unhealthy"
	DATABASECONNECTED databaseStatus = "connected"
)

type healthStatus string

type databaseStatus string

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {

	var health Health

	defer func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, "Failed to encode error", http.StatusInternalServerError)
		}
	}(w, r)

	w.Header().Set("Content-Type", "application/json")

	if err := s.repo.DB.Ping(); err != nil {
		health.Status = STATUSUNHEALTHY
		health.Err = err.Error()
		health.Timestamp = time.Now().Format(time.RFC3339)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	version, dirty, err := repository.GetMigrationInfo(s.repo.DSN)
	if err != nil {
		health.Status = STATUSUNHEALTHY
		health.Err = fmt.Sprintf("migration check failed: %v", err)
		health.Timestamp = time.Now().Format(time.RFC3339)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	health.Status = STATUSHEALTHY
	health.Database = DATABASECONNECTED
	health.Migration_version = version
	health.Migration_dirty = dirty
	health.Timestamp = time.Now().Format(time.RFC3339)
}
