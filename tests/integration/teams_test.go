package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RomanKovalev007/pull_request_service/include/models"
)

func TestCreateTeam_Success(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	team := models.Team{
		TeamName: "test-team-success",
		Members: []models.TeamMember{
			{UserID: "test-user-success-1", Username: "Test User Success 1", IsActive: true},
			{UserID: "test-user-success-2", Username: "Test User Success 2", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	req := httptest.NewRequest("POST", "/team/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}

	var response struct {
		Team models.Team `json:"team"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Team.TeamName != team.TeamName {
		t.Errorf("Expected team name %s, got %s", team.TeamName, response.Team.TeamName)
	}

	if len(response.Team.Members) != len(team.Members) {
		t.Errorf("Expected %d members, got %d", len(team.Members), len(response.Team.Members))
	}

	t.Logf("Team created successfully: %s with %d members", response.Team.TeamName, len(response.Team.Members))
}

func TestCreateTeam_AlreadyExists(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	// Сначала создаем команду
	team := models.Team{
		TeamName: "test-team-duplicate",
		Members: []models.TeamMember{
			{UserID: "test-user-dup-1", Username: "Test User Dup 1", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	req := httptest.NewRequest("POST", "/team/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create initial team: %s", rr.Body.String())
	}

	// Пытаемся создать команду с тем же именем
	req = httptest.NewRequest("POST", "/team/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400 for duplicate team, got %v", status)
	}

	t.Log("Duplicate team creation correctly rejected")
}

func TestGetTeam_Success(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	// Сначала создаем команду
	team := models.Team{
		TeamName: "test-team-get-success",
		Members: []models.TeamMember{
			{UserID: "get-user-success-1", Username: "Get User Success 1", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	req := httptest.NewRequest("POST", "/team/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create team for get test: %s", rr.Body.String())
	}

	// Получаем команду
	req = httptest.NewRequest("GET", "/team/get?team_name=test-team-get-success", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}

	var response models.Team
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.TeamName != team.TeamName {
		t.Errorf("Expected team name %s, got %s", team.TeamName, response.TeamName)
	}

	if len(response.Members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(response.Members))
	}

	t.Logf("Team retrieved successfully: %s", response.TeamName)
}

func TestGetTeam_NotFound(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	req := httptest.NewRequest("GET", "/team/get?team_name=nonexistent-team", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent team, got %v", status)
	}

	t.Log("Non-existent team correctly returns 404")
}
