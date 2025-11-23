package service

import "errors"

var ErrInvalidInput = errors.New("INVALID_INPUT")

type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Code + ": " + e.Message
}
