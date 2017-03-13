package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ValidationErrorFieldErrors map[string]string
type ValidationErrors map[string]ValidationErrorFieldErrors
type APIError struct {
	LowerError error             `json:"-"`
	Type       string            `json:"type"`
	Title      string            `json:"title"`
	Status     int               `json:"status"`
	Detail     string            `json:"detail"`
	ErrorCode  int               `json:"errorCode"`
	Validation *ValidationErrors `json:"validation_messages,omitempty"`
}

func (e APIError) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

func (e APIError) NotFound() bool {
	// {"title":"Not Found","status":404,"detail":"Entity not found."} => endpoint correct but we could not find it
	// {"title":"Not Found","status":404,"detail":"Page not found."} => endpoint incorrect ;(

	return (e.Status == http.StatusNotFound && e.Detail == "Entity not found.")
}

type ClientError struct {
	LowerError error
	Message    string
	Code       int
}

func (e ClientError) Error() string {
	if e.LowerError != nil {
		return fmt.Sprintf("Generic client error [%d:%s] %v", e.Code, e.Message, e.LowerError)
	}
	return fmt.Sprintf("Generic client error [%d:%s]", e.Code, e.Message)
}

type JSONClientError struct {
	LowerError error
	Message    string
	Code       int
}

func (e JSONClientError) Error() string {
	if e.LowerError != nil {
		return fmt.Sprintf("Invalid json recieved [%d:%s] %v", e.Code, e.Message, e.LowerError)
	}
	return fmt.Sprintf("Invalid json recieved [%d:%s]", e.Code, e.Message)
}
