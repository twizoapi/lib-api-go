package twizo

import (
	"encoding/json"
	"net/http"
)

// ApplicationVerifyCredentialsRequest empty struct placeholder (future use)
type ApplicationVerifyCredentialsRequest struct{}

// Submit wil submit the balance request
func (request *ApplicationVerifyCredentialsRequest) Submit() (*ApplicationVerifyCredentialsResponse, error) {
	response := &ApplicationVerifyCredentialsResponse{}
	response.isKeyValid = true

	apiURL, err := GetURLFor("application/verifycredentials")
	if err != nil {
		return nil, err
	}

	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		apiURL,
		nil,
		http.StatusOK,
		response,
	)

	if apiError, ok := err.(*APIError); ok {
		if apiError.NotAuthorized() {
			response.isKeyValid = false
			err = nil
		}
	}

	return response, err
}

type jsonApplicationVerifyCredentialsResponse struct {
	IsTestKey      bool   `json:"isTestKey"`
	ApplicationTag string `json:"applicationTag"`
}

// ApplicationVerifyCredentialsResponse struct that the server returns for a verification request
type ApplicationVerifyCredentialsResponse struct {
	isTestKey      bool
	applicationTag string
	isKeyValid     bool
}

// UnmarshalJSON unmarshals the json returned by the server into the BalanceGetResponse structs
func (response *ApplicationVerifyCredentialsResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonApplicationVerifyCredentialsResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *ApplicationVerifyCredentialsResponse) copyFrom(j *jsonApplicationVerifyCredentialsResponse) error {
	var err error // default err is nil

	response.applicationTag = j.ApplicationTag
	response.isTestKey = j.IsTestKey

	return err
}

// IsTestKey is the current key used a test key
func (response ApplicationVerifyCredentialsResponse) IsTestKey() bool {
	return response.isTestKey
}

// GetApplicationTag get the current credit currency code
func (response ApplicationVerifyCredentialsResponse) GetApplicationTag() string {
	return response.applicationTag
}

// IsKeyValid is the current key valid
func (response ApplicationVerifyCredentialsResponse) IsKeyValid() bool {
	return response.isKeyValid
}

// NewApplicationVerifyCredentials creates a new ApplicationVerifyCredentials
func NewApplicationVerifyCredentials() *ApplicationVerifyCredentialsRequest {
	request := &ApplicationVerifyCredentialsRequest{}
	return request
}

// ApplicationVerifyCredentials retrieves the credit balance of the api key
func ApplicationVerifyCredentials() (*ApplicationVerifyCredentialsResponse, error) {
	request := NewApplicationVerifyCredentials()
	return request.Submit()
}
