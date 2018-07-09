package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type jsonTotpRequest struct {
	Identifier string `json:"identifier"`
	Issuer     string `json:"issuer"`
}

func (request *jsonTotpRequest) copyFrom(r *TotpRequest) {
	request.Identifier = r.identifier
	request.Issuer = r.issuer
}

// TotpRequest request for creating backup codes for id
type TotpRequest struct {
	identifier string
	issuer     string
}

// MarshalJSON is used to convert TotpRequest to json
// MarshalJSON is used to convert NumberLookupRequest to json
func (request *TotpRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := &jsonTotpRequest{}
	jsonRequest.copyFrom(request)
	return json.Marshal(jsonRequest)
}

// GetIdentifier return the identifier of the request
func (request TotpRequest) GetIdentifier() string {
	return request.identifier
}

// GetIssuer return the identifier of the request
func (request TotpRequest) GetIssuer() string {
	return request.issuer
}

// SetIssuer sets the issuer
func (request *TotpRequest) SetIssuer(issuer string) {
	request.issuer = issuer

}

// Create a totpCreate code for an identifier
func (request *TotpRequest) Create(issuer string) (*TotpResponse, error) {
	request.SetIssuer(issuer)

	response := &TotpResponse{}

	apiURL, err := GetURLFor("totp")
	if err != nil {
		return nil, err
	}

	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodPost,
		apiURL,
		request,
		http.StatusCreated,
		response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Check retrieves some information about the id (previousely created with Create)
func (request TotpRequest) Check() (*TotpResponse, error) {
	response := &TotpResponse{}
	apiURL, err := GetURLFor(fmt.Sprintf("totp/%s", url.PathEscape(request.GetIdentifier())))
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

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Verify token for an identifier
func (request TotpRequest) Verify(token string) (*TotpResponse, error) {
	response := &TotpResponse{}

	apiURL, err := GetURLFor(fmt.Sprintf("totp/%s", url.PathEscape(request.GetIdentifier())))
	if err != nil {
		return nil, err
	}

	q := apiURL.Query()
	q.Set("token", token)
	apiURL.RawQuery = q.Encode()

	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		apiURL,
		nil,
		http.StatusOK,
		response,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Delete the existing (if any) backup codes for identifier
func (request TotpRequest) Delete() error {
	apiURL, err := GetURLFor(fmt.Sprintf("totp/%s", url.PathEscape(request.GetIdentifier())))
	if err != nil {
		return err
	}

	return GetClient(RegionCurrent, APIKey).Call(
		http.MethodDelete,
		apiURL,
		request,
		http.StatusNoContent,
		nil,
	)
}

type jsonTotpResponseEmbedded struct {
	JSONVerificationResponse *jsonVerificationResponse `json:"verification,omitempty"`
}
type jsonTotpResponse struct {
	Identifier string                    `json:"identifier"`
	Issuer     string                    `json:"issuer"`
	URL        string                    `json:"uri"`
	Embedded   *jsonTotpResponseEmbedded `json:"_embedded,omitempty"`
	Links      HATEOASLinks              `json:"_links"`
}

// TotpResponse ...
type TotpResponse struct {
	identifier           string
	issuer               string
	url                  *url.URL
	verificationResponse *VerificationResponse
	links                HATEOASLinks
}

// UnmarshalJSON the json response to struct
func (response *TotpResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonTotpResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *TotpResponse) copyFrom(j *jsonTotpResponse) error {

	response.identifier = j.Identifier
	response.issuer = j.Issuer
	response.links = j.Links
	if len(j.URL) > 0 {
		u, err := url.Parse(j.URL)
		if err != nil {
			return err
		}
		response.url = u
	}

	response.verificationResponse = &VerificationResponse{}
	if j.Embedded != nil && j.Embedded.JSONVerificationResponse != nil {
		return response.verificationResponse.copyFrom(j.Embedded.JSONVerificationResponse)
	}

	return nil
}

// GetURL returns the uri of the request
func (response TotpResponse) GetURL() *url.URL {
	return response.url
}

// GetURLSecret returns the secret of the totp URL generated
func (response TotpResponse) GetURLSecret() *string {
	if response.url == nil {
		return nil
	}
	secret := response.url.Query().Get("secret")
	return &secret
}

// GetIssuer returns the issuer of the response
func (response TotpResponse) GetIssuer() string {
	return response.issuer
}

// GetIdentifier returns the identifier of the response
func (response TotpResponse) GetIdentifier() string {
	return response.identifier
}

// GetVerificationResponse returns the verification response or nil
func (response TotpResponse) GetVerificationResponse() *VerificationResponse {
	return response.verificationResponse
}

// NewTotpRequest creates a new BackupCodeRequest
func NewTotpRequest(id string) *TotpRequest {
	request := &TotpRequest{
		identifier: id,
	}
	return request
}

// TotpCreate creates new totpCreate for an identifier
func TotpCreate(id string, issuer string) (*TotpResponse, error) {
	request := NewTotpRequest(id)
	return request.Create(issuer)
}

// TotpCheck creates new totpCreate for an identifier
func TotpCheck(id string) (*TotpResponse, error) {
	request := NewTotpRequest(id)
	return request.Check()
}

// TotpDelete will delete the totpCreate for the identifier supplied
func TotpDelete(id string) error {
	request := NewTotpRequest(id)
	return request.Delete()
}

// TotpVerify will verify a token for an totpCreate
func TotpVerify(id string, token string) (*TotpResponse, error) {
	request := NewTotpRequest(id)
	return request.Verify(token)
}
