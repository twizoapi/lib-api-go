package twizo

import (
	"encoding/json"
	"net/http"
	"time"
)

// RegistrationWidgetSessionRequest contains the struct of the request send for a verification
type RegistrationWidgetSessionRequest struct {
	recipient            Recipient
	allowedTypes         VerificationTypes
	backupCodeIdentifier string
	totpIdentifier       string
	issuer               string
}

type jsonRegistrationWidgetSessionRequest struct {
	Recipient            Recipient         `json:"recipient,omitempty"`
	AllowedTypes         VerificationTypes `json:"allowedTypes,omitempty"`
	BackupCodeIdentifier string            `json:"backupCodeIdentifier,omitempty"`
	TotpIdentifier       string            `json:"totpIdentifier,omitempty"`
	Issuer               string            `json:"issuer,omitempty"`
}

// MarshalJSON is used to convert SmsRequest to json
func (request *RegistrationWidgetSessionRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonRegistrationWidgetSessionRequest{
		Recipient:    request.recipient,
		AllowedTypes: request.allowedTypes,
		Issuer:       request.issuer,
	}

	jsonRequest.BackupCodeIdentifier = request.backupCodeIdentifier
	jsonRequest.TotpIdentifier = request.totpIdentifier

	return json.Marshal(jsonRequest)
}

// SetRecipient set the recipient of the verification
func (request *RegistrationWidgetSessionRequest) SetRecipient(recipient Recipient) {
	request.recipient = recipient
}

// GetRecipient get the recipient of the verification
func (request RegistrationWidgetSessionRequest) GetRecipient() Recipient {
	return request.recipient
}

// GetBackupCodeIdentifier gets the backup code identifier
func (request RegistrationWidgetSessionRequest) GetBackupCodeIdentifier() string {
	return request.backupCodeIdentifier
}

// SetBackupCodeIdentifier sets the backup code identifier
func (request *RegistrationWidgetSessionRequest) SetBackupCodeIdentifier(backupCodeIdentifier string) {
	request.backupCodeIdentifier = backupCodeIdentifier
}

// GetTotpIdentifier gets the totpCreate identifier
func (request RegistrationWidgetSessionRequest) GetTotpIdentifier() string {
	return request.totpIdentifier
}

// SetTotpIdentifier sets the backup code identifier
func (request *RegistrationWidgetSessionRequest) SetTotpIdentifier(totpIdentifier string) {
	request.totpIdentifier = totpIdentifier
}

// GetAllowedTypes get the allowed types
func (request *RegistrationWidgetSessionRequest) GetAllowedTypes() VerificationTypes {
	return request.allowedTypes
}

// SetAllowedTypes is a shorthand to set the allowedTypes
func (request *RegistrationWidgetSessionRequest) SetAllowedTypes(verifications interface{}) {
	switch tVerifications := verifications.(type) {
	case VerificationTypes:
		request.allowedTypes = tVerifications
	case []string:
		for _, element := range tVerifications {
			request.allowedTypes.Add(VerificationType(element))
		}
	}
}

// GetIssuer returns the issuer string
func (request RegistrationWidgetSessionRequest) GetIssuer() string {
	return request.issuer
}

// SetIssuer sets the issuer string for line, push and telegram
func (request *RegistrationWidgetSessionRequest) SetIssuer(issuer string) {
	request.issuer = issuer
}

// Submit wil submit the verification request
func (request *RegistrationWidgetSessionRequest) Submit() (*RegistrationWidgetSessionResponse, error) {
	response := &RegistrationWidgetSessionResponse{}
	response.allowedTypes = &VerificationTypes{}
	response.requestedTypes = &VerificationTypes{}

	apiURL, err := GetURLFor("widget-register-verification/session")
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

// String converts the verification request to a string
func (request RegistrationWidgetSessionRequest) String() string {
	ret, _ := json.Marshal(request)
	return string(ret)
}

// RegistrationWidgetSessionResponse struct that the server returns for a verification request
type RegistrationWidgetSessionResponse struct {
	sessionToken         string
	applicationTag       string
	createdDateTime      time.Time
	issuer               string
	language             string
	recipient            Recipient
	requestedTypes       *VerificationTypes
	registeredTypes      *VerificationTypes
	allowedTypes         *VerificationTypes
	statusMsg            string
	statusCode           VerificationStatusCode
	backupCodeIdentifier string
	totpIdentifier       string
	links                HATEOASLinks
}

type jsonRegistrationWidgetSessionResponse struct {
	SessionToken         string                 `json:"sessionToken"`
	ApplicationTag       string                 `json:"applicationTag"`
	CreatedDateTime      time.Time              `json:"createdDateTime"`
	Issuer               string                 `json:"issuer,omitempty"`
	Language             string                 `json:"language,omitempty"`
	Recipient            Recipient              `json:"recipient"`
	RequestedTypes       *VerificationTypes     `json:"requestedTypes,omitempty"`
	AllowedTypes         *VerificationTypes     `json:"allowedTypes"`
	RegisteredTypes      *VerificationTypes     `json:"registeredTypes"`
	StatusMsg            string                 `json:"status"`
	StatusCode           VerificationStatusCode `json:"statusCode"`
	BackupCodeIdentifier string                 `json:"backupCodeIdentifier,omitempty"`
	TotpIdentifier       string                 `json:"totpIdentifier,omitempty"`
	Links                HATEOASLinks           `json:"_links"`
}

// UnmarshalJSON unmarshals the json returned by the server into the VerificationResponse struct
func (response *RegistrationWidgetSessionResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonRegistrationWidgetSessionResponse{}
	jsonResponse.AllowedTypes = &VerificationTypes{}
	jsonResponse.RequestedTypes = &VerificationTypes{}
	jsonResponse.RegisteredTypes = &VerificationTypes{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *RegistrationWidgetSessionResponse) copyFrom(j *jsonRegistrationWidgetSessionResponse) error {
	var err error // default err is nil

	response.sessionToken = j.SessionToken
	response.applicationTag = j.ApplicationTag
	response.createdDateTime = j.CreatedDateTime
	response.backupCodeIdentifier = j.BackupCodeIdentifier
	response.totpIdentifier = j.TotpIdentifier
	response.issuer = j.Issuer
	response.language = j.Language
	response.recipient = j.Recipient
	response.statusMsg = j.StatusMsg
	response.statusCode = j.StatusCode
	response.links = j.Links
	return err
}

// GetSessionToken gets the session token
func (response RegistrationWidgetSessionResponse) GetSessionToken() string {
	return response.sessionToken
}

// GetApplicationTag gets the application tag set on the verification request
func (response RegistrationWidgetSessionResponse) GetApplicationTag() string {
	return response.applicationTag
}

// GetCreateDateTime gets the time.Time that the request was created
func (response RegistrationWidgetSessionResponse) GetCreateDateTime() time.Time {
	return response.createdDateTime
}

// GetLanguage gets the language that the verification request was made in (not implemented yet)
func (response RegistrationWidgetSessionResponse) GetLanguage() string {
	return response.language
}

// GetRecipient get the recipient of the verification
func (response RegistrationWidgetSessionResponse) GetRecipient() Recipient {
	return response.recipient
}

// GetStatusMsg get the status message
func (response RegistrationWidgetSessionResponse) GetStatusMsg() string {
	return response.statusMsg
}

// GetStatusCode get the status code
func (response RegistrationWidgetSessionResponse) GetStatusCode() VerificationStatusCode {
	return response.statusCode
}

// GetBackupCodeIdentifier gets the backupCodeIdentifier
func (response RegistrationWidgetSessionResponse) GetBackupCodeIdentifier() string {
	return response.backupCodeIdentifier
}

// String stringify a validation response
func (response RegistrationWidgetSessionResponse) String() string {
	ret, _ := json.Marshal(response)
	return string(ret)
}

// NewRegistrationWidgetSessionRequest creates a new widgetsession using a recipient (the only required var)
func NewRegistrationWidgetSessionRequest() *RegistrationWidgetSessionRequest {
	registrationWidgetSessionRequest := &RegistrationWidgetSessionRequest{}
	registrationWidgetSessionRequest.allowedTypes = VerificationTypes{}
	return registrationWidgetSessionRequest
}
