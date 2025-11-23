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

type Service interface {
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
	defaultIdleTimeout   = time.Second * 30
)

type Server struct {
	srv  *http.Server
	repo *repository.Repo
	mux  *http.ServeMux

	teamService  *service.TeamService
	userService  *service.UserService
	prService    *service.PrService
	statsService *service.StatsService

	// Добавляем поля для обработчиков
	teamHandler  *TeamHandler
	userHandler  *UserHandler
	prHandler    *PRHandler
	statsHandler *StatsHandler
}

func NewServer(port string, db *repository.Repo) *Server {
	mux := http.NewServeMux()

	srv := http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultHeaderTimeout,
	}

	server := &Server{
		srv:          &srv,
		repo:         db,
		mux:          mux,
		teamService:  service.NewTeamService(db.TeamRepository),
		userService:  service.NewUserService(db.UserRepository),
		prService:    service.NewPrService(db.PrRepository),
		statsService: service.NewStatsService(db.StatsRepository),
	}

	server.teamHandler = NewTeamHandler(server.teamService)
	server.userHandler = NewUserHandler(server.userService)
	server.prHandler = NewPRHandler(server.prService)
	server.statsHandler = NewStatsHandler(server.statsService)

	return server
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) GetRouter() http.Handler {
	return s.mux
}

func (s *Server) RegisterHandlers() error {
	if s.mux == nil {
		s.mux = http.NewServeMux()
	}

	s.mux.HandleFunc("/health", s.HealthCheck)
	s.mux.HandleFunc("/stats", s.statsHandler.GetStats)

	s.mux.HandleFunc("/team/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.teamHandler.AddTeam(w, r)
	})

	s.mux.HandleFunc("/team/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.teamHandler.GetTeam(w, r)
	})

	s.mux.HandleFunc("/users/setIsActive", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.userHandler.SetUserIsActive(w, r)
	})

	s.mux.HandleFunc("/users/getReview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.userHandler.GetUserPullRequests(w, r)
	})

	s.mux.HandleFunc("/pullRequest/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.prHandler.CreatePullRequest(w, r)
	})

	s.mux.HandleFunc("/pullRequest/merge", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.prHandler.MergePullRequest(w, r)
	})

	s.mux.HandleFunc("/pullRequest/reassign", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.prHandler.ReassignReviewer(w, r)
	})

	s.srv.Handler = s.mux
	return nil
}
