package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

func TestSetUserActive_Deactivate(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}
	team := models.Team{
		TeamName: "active-test-team",
		Members: []models.TeamMember{
			{UserID: "active-user-1", Username: "Active User 1", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	req := httptest.NewRequest("POST", "/team/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create team: %s", rr.Body.String())
	}

	deactivateReq := transport.UserSetActiveRequest{
		UserID:   "active-user-1",
		IsActive: false,
	}

	body, _ = json.Marshal(deactivateReq)
	req = httptest.NewRequest("POST", "/users/setIsActive", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}

	var response struct {
		User models.User `json:"user"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.User.IsActive {
		t.Error("Expected user to be deactivated")
	}

	t.Logf("User deactivated successfully: %s", response.User.UserID)
}

func TestSetUserActive_Reactivate(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	reactivateReq := transport.UserSetActiveRequest{
		UserID:   "active-user-1",
		IsActive: true,
	}

	body, _ := json.Marshal(reactivateReq)
	req := httptest.NewRequest("POST", "/users/setIsActive", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	var response struct {
		User models.User `json:"user"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.User.IsActive {
		t.Error("Expected user to be activated")
	}

	t.Logf("User reactivated successfully: %s", response.User.UserID)
}

func TestSetUserActive_NotFound(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	deactivateReq := transport.UserSetActiveRequest{
		UserID:   "nonexistent-user",
		IsActive: false,
	}

	body, _ := json.Marshal(deactivateReq)
	req := httptest.NewRequest("POST", "/users/setIsActive", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent user, got %v", status)
	}

	t.Log("Non-existent user correctly returns 404")
}

func TestGetUserPullRequests_Success(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	req := httptest.NewRequest("GET", "/users/getReview?user_id=user2", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}

	var response transport.UserPRsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.UserID != "user2" {
		t.Errorf("Expected user ID user2, got %s", response.UserID)
	}

	t.Logf("User PRs retrieved successfully: %s has %d PRs", response.UserID, len(response.PullRequests))
}

func TestGetUserPullRequests_NotFound(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	req := httptest.NewRequest("GET", "/users/getReview?user_id=nonexistent-user", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent user, got %v", status)
	}

	t.Log("Non-existent user PRs correctly returns 404")
}
