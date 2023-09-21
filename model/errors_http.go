package model

import "net/http"

var (
	ErrUnhealthy = NewError(http.StatusInternalServerError, "something went wrong")

	ErrUserWithThisNameExist = NewError(http.StatusBadRequest, "User with this username exist")
	ErrUserNotExist          = NewError(http.StatusBadRequest, "User not exist")
)

type Error interface {
	error
}

type StatusError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (se StatusError) Error() string {
	return se.Message
}

func NewError(code int, message string) Error {
	return StatusError{
		Code:    code,
		Message: message,
	}
}
