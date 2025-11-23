package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

func TestCreatePullRequest_Success(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	createReq := transport.CreatePRRequest{
		PullRequestID:   "test-pr-create-1",
		PullRequestName: "Test PullRequest Creation",
		AuthorID:        "user1",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/pullRequest/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}

	var response transport.CreatePRResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.PullRequest.PullRequestID != createReq.PullRequestID {
		t.Errorf("Expected PullRequest ID %s, got %s", createReq.PullRequestID, response.PullRequest.PullRequestID)
	}

	if response.PullRequest.Status != "OPEN" {
		t.Errorf("Expected status OPEN, got %s", response.PullRequest.Status)
	}

	if len(response.PullRequest.AssignedReviewers) == 0 {
		t.Error("Expected at least one assigned reviewer")
	}

	t.Logf("PullRequestPullRequest created successfully: %s with %d reviewers", response.PullRequest.PullRequestID, len(response.PullRequest.AssignedReviewers))
}

func TestCreatePullRequest_Duplicate(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	createReq := transport.CreatePRRequest{
		PullRequestID:   "test-pr-duplicate",
		PullRequestName: "Test PullRequest Duplicate",
		AuthorID:        "user1",
	}

	body, _ := json.Marshal(createReq)
	
	req := httptest.NewRequest("POST", "/pullRequest/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Failed to create initial PullRequest: %s", rr.Body.String())
	}

	req = httptest.NewRequest("POST", "/pullRequest/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate PullRequest, got %v", status)
	}

	t.Log("PullRequestDuplicate PullRequest creation correctly rejected")
}

func TestCreatePullRequest_AuthorNotFound(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	createReq := transport.CreatePRRequest{
		PullRequestID:   "test-pr-author-not-found",
		PullRequestName: "Test PullRequest Author Not Found",
		AuthorID:        "nonexistent-author",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/pullRequest/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent author, got %v", status)
	}

	t.Log("PullRequestNon-existent author correctly returns 404")
}

func TestMergePullRequest_Success(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}


	mergeReq := transport.MergePRRequest{
		PullRequestID: "pr6",
	}

	body, _ := json.Marshal(mergeReq)
	req := httptest.NewRequest("POST", "/pullRequest/merge", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}

	var response transport.MergePRResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.PullRequest.Status != "MERGED" {
		t.Errorf("Expected status MERGED, got %s", response.PullRequest.Status)
	}

	if response.PullRequest.MergedAt == nil {
		t.Error("Expected MergedAt to be set")
	}

	t.Logf("PullRequestPullRequest merged successfully: %s at %v", response.PullRequest.PullRequestID, response.PullRequest.MergedAt)
}

func TestMergePullRequest_NotFound(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	mergeReq := transport.MergePRRequest{
		PullRequestID: "nonexistent-pr",
	}

	body, _ := json.Marshal(mergeReq)
	req := httptest.NewRequest("POST", "/pullRequest/merge", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent PullRequest, got %v", status)
	}

	t.Log("PullRequestNon-existent PullRequest merge correctly returns 404")
}

func TestReassignReviewer_Success(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	reassignReq := transport.ReassignRequest{
		PullRequestID: "pr1",
		OldUserID:     "user2",
	}

	body, _ := json.Marshal(reassignReq)
	req := httptest.NewRequest("POST", "/pullRequest/reassign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}

	var response transport.ReassignResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.ReplacedBy == "" {
		t.Error("Expected new reviewer ID")
	}

	if response.ReplacedBy == reassignReq.OldUserID{
		t.Error("Expected different reviewer after reassignment")
	}

	t.Logf("PullRequestReviewer reassigned successfully: %s -> %s", reassignReq.OldUserID, response.ReplacedBy)
}

func TestReassignReviewer_NotFound(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	reassignReq := transport.ReassignRequest{
		PullRequestID: "nonexistent-pr",
		OldUserID:     "user1",
	}

	body, _ := json.Marshal(reassignReq)
	req := httptest.NewRequest("POST", "/pullRequest/reassign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent PullRequest, got %v", status)
	}

	t.Log("PullRequestNon-existent PullRequest reassign correctly returns 404")
}

func TestReassignReviewer_NotAssigned(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	reassignReq := transport.ReassignRequest{
		PullRequestID: "pr5",
		OldUserID:     "user7",
	}

	body, _ := json.Marshal(reassignReq)
	req := httptest.NewRequest("POST", "/pullRequest/reassign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Expected status 409 for not assigned reviewer, got %v", status)
	}

	t.Log("PullRequestNot assigned reviewer reassign correctly returns 409")
}