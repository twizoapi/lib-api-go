package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type jsonBackupCodeRequest struct {
	Identifier string `json:"identifier"`
}

// BackupCodeRequest request for creating backup codes for id
type BackupCodeRequest struct {
	identifier string
}

// MarshalJSON is used to convert NumberLookupRequest to json
func (request *BackupCodeRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonBackupCodeRequest{
		Identifier: request.identifier,
	}

	return json.Marshal(jsonRequest)
}

// Delete the existing (if any) backup codes for identifier
func (request *BackupCodeRequest) Delete() error {
	apiURL, err := GetURLFor(fmt.Sprintf("backupcode/%s", url.PathEscape(request.GetIdentifier())))
	if err != nil {
		return err
	}

	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodDelete,
		apiURL,
		request,
		http.StatusNoContent,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// Status will retrieve the backupcode status
func (request *BackupCodeRequest) Status() (*BackupCodeResponse, error) {
	response := &BackupCodeResponse{}
	apiURL, err := GetURLFor(fmt.Sprintf("backupcode/%s", url.PathEscape(request.GetIdentifier())))
	if err != nil {
		return nil, err
	}

	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		apiURL,
		request,
		http.StatusOK,
		response,
	)

	return response, err
}

// Verify token did belong to the identifier of the request
func (request *BackupCodeRequest) Verify(token string) (*BackupCodeResponse, error) {
	response := &BackupCodeResponse{}

	apiURL, err := GetURLFor(fmt.Sprintf("backupcode/%s", url.PathEscape(request.GetIdentifier())))
	if err != nil {
		return nil, err
	}

	q := apiURL.Query()
	q.Set("token", token)
	apiURL.RawQuery = q.Encode()

	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		apiURL,
		request,
		http.StatusOK,
		response,
	)

	if apiError, ok := err.(*APIError); ok {
		response.verificationResponse = &VerificationResponse{}
		if apiError.Status() == http.StatusNotFound {
			response.verificationResponse.statusCode = VerificationTokenInvalid
		} else if apiError.Status() == http.StatusUnprocessableEntity &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenInvalid {
			response.verificationResponse.statusCode = VerificationTokenInvalid
		} else if apiError.Status() == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenExpired {
			response.verificationResponse.statusCode = VerificationTokenExpired
		} else if apiError.Status() == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenAlreadyVerified {
			response.verificationResponse.statusCode = VerificationTokenAlreadyVerified
		} else if apiError.Status() == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenFailed {
			response.verificationResponse.statusCode = VerificationTokenFailed
		} else {
			// undocumented response, error out
			return nil, err
		}
		return response, nil
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Update the existing (or create new ones if they do not exist) backup tokens
func (request *BackupCodeRequest) Update() (*BackupCodeResponse, error) {
	return request.doUpdateCreate(true)
}

// Create wil submit the balance request
func (request *BackupCodeRequest) Create() (*BackupCodeResponse, error) {
	return request.doUpdateCreate(false)
}

func (request *BackupCodeRequest) doUpdateCreate(update bool) (*BackupCodeResponse, error) {
	response := &BackupCodeResponse{}

	urlPart := "backupcode"
	method := http.MethodPost
	expect := http.StatusCreated
	if update {
		urlPart = fmt.Sprintf(
			"backupcode/%s",
			url.PathEscape(request.GetIdentifier()),
		)
		method = http.MethodPut
		expect = http.StatusOK
	}

	apiURL, err := GetURLFor(urlPart)
	if err != nil {
		return nil, err
	}

	err = GetClient(RegionCurrent, APIKey).Call(
		method,
		apiURL,
		request,
		expect,
		response,
	)

	if apiError, ok := err.(*APIError); !update && ok && apiError.Conflict() {
		response.alreadyExists = true
		response.identifier = request.GetIdentifier()
		err = nil
	} else {
		response.alreadyExists = false
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetIdentifier returns the identifier of the current backup code request
func (request *BackupCodeRequest) GetIdentifier() string {
	return request.identifier
}

type jsonBackupCodeResponseEmbedded struct {
	JSONVerificationResponse *jsonVerificationResponse `json:"verification,omitempty"`
}

type jsonBackupCodeResponse struct {
	Identifier        string                         `json:"identifier"`
	AmountOfCodesLeft int                            `json:"amountOfCodesLeft"`
	Codes             []string                       `json:"codes"`
	CreateDateTime    *time.Time                     `json:"createdDateTime,omitempty"`
	Embedded          jsonBackupCodeResponseEmbedded `json:"_embedded"`
	Links             HATEOASLinks                   `json:"_links"`
}

// BackupCodeResponse response for backup create api call
type BackupCodeResponse struct {
	identifier           string
	amountOfCodesLeft    int
	codes                []string
	createDateTime       *time.Time
	links                HATEOASLinks
	verificationResponse *VerificationResponse
	alreadyExists        bool
}

// UnmarshalJSON the json response to struct
func (response *BackupCodeResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonBackupCodeResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *BackupCodeResponse) copyFrom(j *jsonBackupCodeResponse) error {
	var err error

	response.identifier = j.Identifier
	response.codes = j.Codes
	response.amountOfCodesLeft = j.AmountOfCodesLeft
	response.createDateTime = j.CreateDateTime
	response.links = j.Links

	if j.Embedded.JSONVerificationResponse != nil {
		response.verificationResponse = &VerificationResponse{}
		err = response.verificationResponse.copyFrom(j.Embedded.JSONVerificationResponse)
	}

	return err
}

// GetIdentifier gets the identifier sent with the request
func (response BackupCodeResponse) GetIdentifier() string {
	return response.identifier
}

// GetCodes returns the codes that were assigned (only for a create operation)
func (response BackupCodeResponse) GetCodes() []string {
	return response.codes
}

// GetAmountOfCodesLeft returns the amount of codes left to try
func (response BackupCodeResponse) GetAmountOfCodesLeft() int {
	return response.amountOfCodesLeft
}

// GetCreateDateTime returns the date timestamp when the codes were generated
func (response BackupCodeResponse) GetCreateDateTime() *time.Time {
	return response.createDateTime
}

// GetVerificationResponse returns the verification response or nil
func (response *BackupCodeResponse) GetVerificationResponse() *VerificationResponse {
	return response.verificationResponse
}

// AlreadyExists returns false if they already exist
func (response BackupCodeResponse) AlreadyExists() bool {
	return response.alreadyExists
}

// NewBackupCodeRequest creates a new BackupCodeRequest
func NewBackupCodeRequest(id string) *BackupCodeRequest {
	request := &BackupCodeRequest{identifier: id}
	return request
}

// BackupCodeCreate creates new backup codes for an identifier
func BackupCodeCreate(id string) (*BackupCodeResponse, error) {
	request := NewBackupCodeRequest(id)
	return request.Create()
}

// BackupCodeUpdate updates the backup codes for an identifier (this will
// invalidate the old backup codes)
func BackupCodeUpdate(id string) (*BackupCodeResponse, error) {
	request := NewBackupCodeRequest(id)
	return request.Update()
}

// BackupCodeDelete will delete the backup codes for the identifier supplied
func BackupCodeDelete(id string) error {
	request := NewBackupCodeRequest(id)
	return request.Delete()
}

// BackupCodeVerify will verify a token for an intentifier
func BackupCodeVerify(id string, token string) (*BackupCodeResponse, error) {
	request := NewBackupCodeRequest(id)
	return request.Verify(token)
}

// BackupCodeStatus will return backup status
func BackupCodeStatus(id string) (*BackupCodeResponse, error) {
	request := NewBackupCodeRequest(id)
	return request.Status()
}

// BackupCodeAmountLeft returns the amount of codes left or 0 on error
func BackupCodeAmountLeft(id string) (int, error) {
	response, err := BackupCodeStatus(id)
	if err != nil {
		return 0, err
	}
	return response.GetAmountOfCodesLeft(), err
}
