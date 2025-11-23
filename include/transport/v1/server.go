package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	"github.com/RomanKovalev007/pull_request_service/include/repository"
	"github.com/RomanKovalev007/pull_request_service/include/service"
	transport "github.com/RomanKovalev007/pull_request_service/include/transport/models"
)

type Service interface{
	CreateTeam(ctx context.Context, req models.Team) (*transport.TeamCreateResponse, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)

	SetUserIsActive(ctx context.Context, req transport.UserSetActiveRequest) (*transport.UserSetActiveResponse, error)
	GetUserPullRequests(ctx context.Context, userID string) (*transport.UserPRsResponse, error)
	
	CreatePullRequest(ctx context.Context, req transport.CreatePRRequest) (*transport.CreatePRResponse, error) 
	MergePullRequest(ctx context.Context, req transport.MergePRRequest) (*transport.MergePRResponse, error)
	ReassignReviewer(ctx context.Context, req transport.ReassignRequest) (*transport.ReassignResponse, error) 
}


var (
	defaultHeaderTimeout = time.Second * 5
	defaultIdleTimeout = time.Second * 30
)

type Server struct {
	srv *http.Server
	repo  *repository.Repo
}

func NewServer(port string, db *repository.Repo) *Server {
	srv := http.Server{
		Addr:              ":" + port,
		Handler:           nil,
		IdleTimeout: defaultIdleTimeout,
		ReadHeaderTimeout: defaultHeaderTimeout,
	}
	return &Server{
		srv: &srv,
		repo:  db,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) RegisterHandlers() error {
	statsRepo := repository.NewStatsRepository(s.repo.DB)
    statsService := service.NewStatsService(statsRepo)
    statsHandler := NewStatsHandler(statsService)

	service := service.NewService(s.repo)

	teamHandler := NewTeamHandler(service)
	userHandler := NewUserHandler(service)
	prHandler := NewPRHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.HealthCheck)

	mux.HandleFunc("/stats", statsHandler.GetStats)

	mux.HandleFunc("/team/add",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		teamHandler.AddTeam(w, r)
	})

	mux.HandleFunc("/team/get",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		teamHandler.GetTeam(w, r)
	})

	mux.HandleFunc("/users/setIsActive",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userHandler.SetUserIsActive(w, r)
	})

	mux.HandleFunc("/users/getReview",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userHandler.GetUserPullRequests(w, r)
	})

	mux.HandleFunc("/pullRequest/create",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		prHandler.CreatePullRequest(w, r)
	})

	mux.HandleFunc("/pullRequest/merge",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		prHandler.MergePullRequest(w, r)
	})

	mux.HandleFunc("/pullRequest/reassign",
	func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		prHandler.ReassignReviewer(w, r)
	})

	s.srv.Handler = mux

	return nil
}