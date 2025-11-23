package integration

import (
	"testing"
)

func TestDBTablesExist(t *testing.T) {
	db := GetTestDB()
	if db == nil {
		t.Fatal("TestDB is nil")
	}

	tables := []string{"teams", "users", "pull_requests", "pr_reviewers"}
	for _, table := range tables {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_name = $1
			)`, table).Scan(&exists)

		if err != nil {
			t.Fatalf("Failed to check table %s: %v", table, err)
		}

		if !exists {
			t.Errorf("Table %s does not exist", table)
		} else {
			t.Logf("Table %s exists", table)
		}
	}
}
