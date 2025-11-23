package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStats_Success(t *testing.T) {
	router := GetTestRouter()
	if router == nil {
		t.Fatal("Test router is not initialized")
	}

	req := httptest.NewRequest("GET", "/stats", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
		return
	}


	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse stats response: %v", err)
	}


	expectedFields := []string{"user_stats", "pr_stats", "team_stats", "total_stats"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Expected field %s in stats response", field)
		}
	}

	t.Log("Stats retrieved successfully")
}