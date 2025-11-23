package integration

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/RomanKovalev007/pull_request_service/include/repository"
	v1 "github.com/RomanKovalev007/pull_request_service/include/transport/v1"
	"github.com/RomanKovalev007/pull_request_service/tests/integration/testutils"
	_ "github.com/lib/pq"
)

var (
	TestDB     *sql.DB
	TestServer *v1.Server
	TestConfig *testutils.DBConfig
)

func setup() {
	var err error

	cfg := &testutils.DBConfig{
		Host:     getEnv("TEST_DB_HOST", "localhost"),
		Port:     getEnv("TEST_DB_PORT", "5432"),
		User:     getEnv("TEST_DB_USER", "postgres"),
		Password: getEnv("TEST_DB_PASSWORD", "password"),
		DBName:   getEnv("TEST_DB_NAME", "pr_reviewer_test"),
	}

	TestConfig = cfg

	TestDB, err = testutils.ConnectTestDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := testutils.RunTestMigrations(TestDB); err != nil {
		log.Fatalf("Failed to run test migrations: %v", err)
	}

	if err := testutils.LoadTestData(TestDB); err != nil {
		log.Fatalf("Failed to load test data: %v", err)
	}

	repoCfg := repository.Config{
		DBHost:     cfg.Host,
		DBPort:     cfg.Port,
		DBUser:     cfg.User,
		DBPassword: cfg.Password,
		DBName:     cfg.DBName,
	}

	repo, err := repository.NewDB(repoCfg)
	if err != nil {
		log.Fatalf("Failed to create newdb: %v", err)
	}

	TestServer = v1.NewServer("8080", repo)
	if err := TestServer.RegisterHandlers(); err != nil {
		log.Fatalf("Failed to register handlers: %v", err)
	}
}

func teardown() {
	if TestDB != nil {
		if err := testutils.CleanTestData(TestDB); err != nil {
			log.Printf("Failed to clean test data: %v", err)
		}
		TestDB.Close()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetTestDB() *sql.DB {
	if TestDB == nil {
		log.Println("WARNING: testDB is nil!")
	}
	return TestDB
}

func GetTestRouter() http.Handler {
	return TestServer.GetRouter()
}
