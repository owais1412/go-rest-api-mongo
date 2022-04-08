package models

type SuccessMessage struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
