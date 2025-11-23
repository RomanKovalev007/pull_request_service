package models

type ErrorResponseErrorCode string

const (
	NOCANDIDATE ErrorResponseErrorCode = "NO_CANDIDATE"
	NOTASSIGNED ErrorResponseErrorCode = "NOT_ASSIGNED"
	NOTFOUND    ErrorResponseErrorCode = "NOT_FOUND"
	PREXISTS    ErrorResponseErrorCode = "PR_EXISTS"
	PRMERGED    ErrorResponseErrorCode = "PR_MERGED"
	TEAMEXISTS  ErrorResponseErrorCode = "TEAM_EXISTS"
	INVALID_INPUT ErrorResponseErrorCode = "INVALID_INPUT"
	INTERNAL_ERROR ErrorResponseErrorCode = "INTERNAL_ERROR"
	STATUS_OK ErrorResponseErrorCode = "STATUS_OK"
)

type ErrorResponse struct {
	Error struct {
		Code    ErrorResponseErrorCode `json:"code"`
		Message string                 `json:"message"`
	} `json:"error"`
}