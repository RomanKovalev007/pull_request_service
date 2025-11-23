package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func ConnectTestDB(cfg *DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func RunTestMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS teams (
            team_name VARCHAR(255) PRIMARY KEY,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,

		`CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(255) PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
            team_name VARCHAR(255) NOT NULL,
            is_active BOOLEAN DEFAULT true,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (team_name) REFERENCES teams(team_name) ON DELETE CASCADE
        )`,

		`CREATE TABLE IF NOT EXISTS pull_requests (
            id VARCHAR(255) PRIMARY KEY,
            pull_request_name VARCHAR(255) NOT NULL,
            author_id VARCHAR(255) NOT NULL,
            status VARCHAR(50) DEFAULT 'OPEN',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            merged_at TIMESTAMP NULL,
            FOREIGN KEY (author_id) REFERENCES users(id)
        )`,

		`CREATE TABLE IF NOT EXISTS pr_reviewers (
            pull_request_id VARCHAR(255),
            reviewer_id VARCHAR(255),
            assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (pull_request_id, reviewer_id),
            FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
            FOREIGN KEY (reviewer_id) REFERENCES users(id)
        )`,

		`CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status)`,
		`CREATE INDEX IF NOT EXISTS idx_pr_author ON pull_requests(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer ON pr_reviewers(reviewer_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

func LoadTestData(db *sql.DB) error {
	teams := []string{"backend", "frontend", "payments", "mobile"}
	for _, team := range teams {
		if _, err := db.Exec("INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING", team); err != nil {
			return fmt.Errorf("failed to insert team %s: %w", team, err)
		}
	}

	users := []struct {
		ID       string
		Username string
		Team     string
		Active   bool
	}{
		{"user1", "Oleg", "backend", true},
		{"user2", "Marat", "backend", true},
		{"user3", "Kostya", "frontend", true},
		{"user4", "David", "payments", true},
		{"user5", "Roman", "payments", true},
		{"user6", "Gleb", "mobile", true},
		{"user7", "Sonya", "mobile", false},
		{"user8", "Nikita", "backend", true},
		{"user9", "Alex", "backend", true},
		{"user10", "Misha", "frontend", true},
	}

	for _, user := range users {
		_, err := db.Exec(`
            INSERT INTO users (id, username, team_name, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE SET 
            username = $2, team_name = $3, is_active = $4`,
			user.ID, user.Username, user.Team, user.Active)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.ID, err)
		}
	}

	prs := []struct {
		ID       string
		Name     string
		AuthorID string
		Status   string
	}{
		{"pr1", "Add authentication", "user1", "OPEN"},
		{"pr2", "Fix database connection", "user2", "OPEN"},
		{"pr3", "Update UI components", "user3", "MERGED"},
		{"pr4", "Payment gateway integration", "user4", "OPEN"},
		{"pr5", "Mobile navigation", "user6", "OPEN"},
		{"pr6", "Linter", "user9", "OPEN"},
	}

	for _, pr := range prs {
		_, err := db.Exec(`
            INSERT INTO pull_requests (id, pull_request_name, author_id, status) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE SET 
            pull_request_name = $2, author_id = $3, status = $4`,
			pr.ID, pr.Name, pr.AuthorID, pr.Status)
		if err != nil {
			return fmt.Errorf("failed to insert PR %s: %w", pr.ID, err)
		}
	}

	reviewers := []struct {
		PRID     string
		Reviewer string
	}{
		{"pr1", "user2"}, {"pr1", "user8"},
		{"pr2", "user1"},
		{"pr3", "user10"},
		{"pr4", "user5"},
		{"pr5", "user7"},
		{"pr6", "user8"}, {"pr6", "user1"},
	}

	for _, review := range reviewers {
		_, err := db.Exec(`
            INSERT INTO pr_reviewers (pull_request_id, reviewer_id) 
            VALUES ($1, $2)
            ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING`,
			review.PRID, review.Reviewer)
		if err != nil {
			return fmt.Errorf("failed to assign reviewer %s to PR %s: %w", review.Reviewer, review.PRID, err)
		}
	}

	log.Printf("Loaded test data: %d teams, %d users, %d PRs", len(teams), len(users), len(prs))
	return nil
}

func CleanTestData(db *sql.DB) error {
	tables := []string{
		"pr_reviewers",
		"pull_requests",
		"users",
		"teams",
	}

	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	log.Println("Cleaned all test data")
	return nil
}
