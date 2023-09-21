package model

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}
