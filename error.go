package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type validationErrorFieldErrors map[string]string
type validationErrors map[string]validationErrorFieldErrors

// APIError struct
type APIError struct {
	LowerError error             `json:"-"`
	Type       string            `json:"type"`
	Title      string            `json:"title"`
	Status     int               `json:"status"`
	Detail     string            `json:"detail"`
	ErrorCode  int               `json:"errorCode"`
	Validation *validationErrors `json:"validation_messages,omitempty"`
}

// Error casts to an actual error struct
func (e APIError) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

// NotFound is true if the endpoint sent to was not found
func (e APIError) NotFound() bool {
	// {"title":"Not Found","status":404,"detail":"Entity not found."} => endpoint correct but we could not find it
	// {"title":"Not Found","status":404,"detail":"Page not found."} => endpoint incorrect ;(

	return e.Status == http.StatusNotFound && e.Detail == "Entity not found."
}

// ClientError struct contains the errors returned when communicating with the server
type ClientError struct {
	LowerError error
	Message    string
	Code       int
}

// Error casts to an actual error struct
func (e ClientError) Error() string {
	if e.LowerError != nil {
		return fmt.Sprintf("Generic client error [%d:%s] %v", e.Code, e.Message, e.LowerError)
	}
	return fmt.Sprintf("Generic client error [%d:%s]", e.Code, e.Message)
}

// JSONClientError struct contains the errors related to erroneous json received from the server
type JSONClientError struct {
	LowerError error
	Message    string
	Code       int
}

// Error casts to an actual error struct
func (e JSONClientError) Error() string {
	if e.LowerError != nil {
		return fmt.Sprintf("Invalid json recieved [%d:%s] %v", e.Code, e.Message, e.LowerError)
	}
	return fmt.Sprintf("Invalid json recieved [%d:%s]", e.Code, e.Message)
}
