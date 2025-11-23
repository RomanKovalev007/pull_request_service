package service

import (
	"context"
	"time"

	"github.com/RomanKovalev007/pull_request_service/include/models"
	"github.com/RomanKovalev007/pull_request_service/include/repository"
)

type StatsRepository interface {
	GetPRStats(ctx context.Context) ([]models.PullRequestStat, error)
	GetTeamStats(ctx context.Context) ([]models.TeamStat, error)
	GetUserStats(ctx context.Context) ([]models.UserStat, error)
}

type StatsService struct {
    statsRepo StatsRepository
}

func NewStatsService(statsRepo *repository.StatsRepository) *StatsService {
    return &StatsService{statsRepo: statsRepo}
}

func (s *StatsService) GetStats(ctx context.Context) (*models.StatsResponse, error) {
    userStats, err := s.statsRepo.GetUserStats(ctx)
    if err != nil {
        return nil, err
    }

    prStats, err := s.statsRepo.GetPRStats(ctx)
    if err != nil {
        return nil, err
    }

    teamStats, err := s.statsRepo.GetTeamStats(ctx)
    if err != nil {
        return nil, err
    }

    totalStats := models.TotalStats{
		TotalUsers: len(userStats),
		TotalPRs: len(prStats),
		TotalTeams: len(teamStats),
	}

	for _, team := range teamStats{
		totalStats.TotalActiveReviewers += team.ActiveReviewers
	}
	
	for _, pr := range prStats{
		if pr.Status == "MERGED"{
			totalStats.MergedPRs += 1
		} else{
			totalStats.OpenPRs += 1
		}
	}

    return &models.StatsResponse{
        UserStats:  userStats,
        PRStats:    prStats,
        TeamStats:  teamStats,
        TotalStats: totalStats,
        Timestamp:  time.Now(),
    }, nil
}