package models

type ErrorsArray struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []ErrorsArray `json:"Errors"`
}
