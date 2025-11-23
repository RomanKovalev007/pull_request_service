// main.go
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/RomanKovalev007/pull_request_service/include/config"
	"github.com/RomanKovalev007/pull_request_service/include/repository"
	v1 "github.com/RomanKovalev007/pull_request_service/include/transport/v1"
	"github.com/joho/godotenv"
)

const (
	defaultTimeout = time.Duration(30) * time.Second
)

func main() {
	// Parse config
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env"
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	cfg.DSN = cfg.FormatConnectionString()

	// Wait for DB start
	time.Sleep(3 * time.Second)

	// Run migrations
	log.Println("Applying database migrations...")
	if err := repository.RunMigrations(cfg.Migration_Path, cfg.DSN); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database
	repo, err := repository.NewDB(cfg.Config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.DB.Close()

	log.Println("Checking if tables were created...")
	checkTables(repo)

	// Initialize server
	server := v1.NewServer(cfg.Port, repo)
	err = server.RegisterHandlers()
	if err != nil {
		log.Fatalf("Failed to register handlers: %v", err)
	}

	// Starting server
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Server starting on %s", cfg.BaseURL)
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Gracefull shutdown
	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	log.Println("Shutdown signal received, starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}

	repo.DB.Close()
	log.Println("DB closed")

	wg.Wait()
	log.Println("Server stopped gracefully")
}

func checkTables(db *repository.Repo) {
	tables := []string{"teams", "users", "pull_requests", "pr_reviewers"}
	for _, table := range tables {
		var exists bool
		err := db.DB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)", table).Scan(&exists)
		if err != nil {
			log.Printf("Error checking table %s: %v", table, err)
		} else {
			log.Printf("Table %s exists: %t", table, exists)
		}
	}
}
