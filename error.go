package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonAPIError struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	ErrorCode int    `json:"errorCode,omitempty"`
}

// APIError struct
type APIError struct {
	title      string
	detail     string
	status     int
	lowerError error
	errorType  string
	errorCode  int
}

// NewAPIError creates a new API error based on the title and the status
func NewAPIError(title string, status int) *APIError {
	return &APIError{
		title:  title,
		status: status,
		detail: title,
	}
}

// UnmarshalJSON unmarshal the json returned by the server into the APIErrorResponse struct
func (e *APIError) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonAPIError{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return e.copyFrom(jsonResponse)
}

func (e *APIError) copyFrom(j *jsonAPIError) error {
	e.errorType = j.Type
	e.title = j.Title
	e.status = j.Status
	e.detail = j.Detail

	return nil
}

// Error casts to an actual error struct
func (e APIError) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

// Status returns the status of the error
func (e APIError) Status() int {
	return e.status
}

// ErrorCode returns the error code of the request
func (e APIError) ErrorCode() int {
	return e.errorCode
}

// Detail returns the detail of the error
func (e APIError) Detail() string {
	return e.detail
}

// Title returns the title of the error
func (e APIError) Title() string {
	return e.title
}

// NotFound is true if the endpoint sent to was not found
func (e APIError) NotFound() bool {
	// {"title":"Not Found","status":404,"detail":"Entity not found."} => endpoint correct but we could not find it
	// {"title":"Not Found","status":404,"detail":"Page not found."} => endpoint incorrect ;(

	return e.Status() == http.StatusNotFound && e.Detail() == "Entity not found."
}

// NotAuthorized is true if api key is not valid
func (e APIError) NotAuthorized() bool {
	return e.Status() == http.StatusUnauthorized
}

// Conflict is true if there was a create confict (it already exists)
func (e APIError) Conflict() bool {
	return e.Status() == http.StatusConflict
}

// UnprocessableEntity is true if there was an unprocessable entity (Incorrect validation)
func (e APIError) UnprocessableEntity() bool {
	return e.Status() == http.StatusUnprocessableEntity
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

//
// API Validation Errors
//
// On the API side the returned content of validation_messages does not have one structure, can have multiple
//   the choice was made to either standardize or add parsing code here later on.
//

// ValidationErrors represents the current unstructured json
type ValidationErrors map[string]interface{}

type jsonAPIValidationError struct {
	jsonAPIError
	Validation ValidationErrors `json:"validation_messages"`
}

// APIValidationError struct
type APIValidationError struct {
	APIError
	validation ValidationErrors
}

// UnmarshalJSON unmarshals the json returned by the server into the APIErrorResponse struct
func (e *APIValidationError) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonAPIValidationError{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return e.copyFrom(jsonResponse)
}

func (e *APIValidationError) copyFrom(j *jsonAPIValidationError) error {
	err := e.APIError.copyFrom(&j.jsonAPIError)
	e.validation = j.Validation

	return err
}

// VerificationErrors returns the unstructured json error struct (interface {})
func (e APIValidationError) VerificationErrors() ValidationErrors {
	return e.validation
}
