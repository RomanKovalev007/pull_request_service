package integration

import (
	"testing"
)

func TestServerInitialization(t *testing.T) {
	if TestServer == nil {
		t.Fatal("Test server should not be nil")
	}

	t.Log("Server and router initialized correctly")
}
