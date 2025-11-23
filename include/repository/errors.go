package repository

import "errors"

var (
	ErrNotFound    = errors.New("NOT_FOUND")
	ErrPRExists    = errors.New("PR_EXISTS")
	ErrTeamExists  = errors.New("TEAM_EXISTS")
	ErrPRMerged    = errors.New("PR_MERGED")
	ErrNotAssigned = errors.New("NOT_ASSIGNED")
	ErrNoCandidate = errors.New("NO_CANDIDATE")
)
