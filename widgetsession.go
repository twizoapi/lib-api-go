package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// WidgetSessionRequest contains the struct of the request send for a verification
type WidgetSessionRequest struct {
	recipient            Recipient
	allowedTypes         VerificationTypes
	backupCodeIdentifier string
	totpIdentifier       string
	tokenLength          int
	tokenType            VerificationTokenType
	tag                  string
	issuer               string
	bodyTemplate         string
	sender               string
	senderNpi            int
	senderTon            int
	dcs                  int
}

type jsonWidgetSessionRequest struct {
	Recipient            Recipient             `json:"recipient,omitempty"`
	AllowedTypes         VerificationTypes     `json:"allowedTypes,omitempty"`
	BackupCodeIdentifier string                `json:"backupCodeIdentifier,omitempty"`
	TotpIdentifier       string                `json:"totpIdentifier,omitempty"`
	TokenLength          int                   `json:"tokenLength,omitempty"`
	TokenType            VerificationTokenType `json:"tokenType,omitempty"`
	Tag                  string                `json:"tag,omitempty"`
	Issuer               string                `json:"issuer,omitempty"`
	BodyTemplate         string                `json:"bodyTemplate,omitempty"`
	Sender               string                `json:"sender,omitempty"`
	SenderNpi            int                   `json:"senderNpi,omitempty"`
	SenderTon            int                   `json:"senderTon,omitempty"`
	Dcs                  int                   `json:"dcs,omitempty"`
}

// MarshalJSON is used to convert SmsRequest to json
func (request *WidgetSessionRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonWidgetSessionRequest{
		Recipient:    request.recipient,
		AllowedTypes: request.allowedTypes,
		TokenLength:  request.tokenLength,
		TokenType:    request.tokenType,
		Tag:          request.tag,
		Issuer:       request.issuer,
		BodyTemplate: request.bodyTemplate,
		Sender:       request.sender,
		SenderTon:    request.senderTon,
		SenderNpi:    request.senderNpi,
		Dcs:          request.dcs,
	}

	jsonRequest.BackupCodeIdentifier = request.backupCodeIdentifier
	jsonRequest.TotpIdentifier = request.totpIdentifier

	return json.Marshal(jsonRequest)
}

// SetRecipient set the recipient of the verification
func (request *WidgetSessionRequest) SetRecipient(recipient Recipient) {
	request.recipient = recipient
}

// GetRecipient get the recipient of the verification
func (request WidgetSessionRequest) GetRecipient() Recipient {
	return request.recipient
}

// SetBodyTemplate set the body template %token% will be replaced by the actual token sent
func (request *WidgetSessionRequest) SetBodyTemplate(template string) {
	request.bodyTemplate = template
}

// GetBodyTemplate get the body template
func (request WidgetSessionRequest) GetBodyTemplate() string {
	return request.bodyTemplate
}

// SetDcs sets the dcs of the verification request
func (request *WidgetSessionRequest) SetDcs(dcs int) {
	request.dcs = dcs
}

// GetDcs get the dcs of the verification request
func (request WidgetSessionRequest) GetDcs() int {
	return request.dcs
}

// SetSender set the sender
func (request *WidgetSessionRequest) SetSender(sender string) {
	request.sender = sender
}

// GetSender get the sender
func (request WidgetSessionRequest) GetSender() string {
	return request.sender
}

// SetSenderNpi sets the senderNpi
func (request *WidgetSessionRequest) SetSenderNpi(npi int) {
	request.senderNpi = npi
}

// GetSenderNpi gets the senderNpi
func (request WidgetSessionRequest) GetSenderNpi() int {
	return request.senderNpi
}

// SetSenderTon sets the senderTon
func (request *WidgetSessionRequest) SetSenderTon(ton int) {
	request.senderTon = ton
}

// GetSenderTon get the senderTon
func (request WidgetSessionRequest) GetSenderTon() int {
	return request.senderTon
}

// SetTag set the tag for this verification
func (request *WidgetSessionRequest) SetTag(tag string) {
	request.tag = tag
}

// GetTag gets the tag set for this verification
func (request WidgetSessionRequest) GetTag() string {
	return request.tag
}

// SetTokenLength sets the token length
func (request *WidgetSessionRequest) SetTokenLength(length int) {
	request.tokenLength = length
}

// GetTokenLength gets the token length
func (request WidgetSessionRequest) GetTokenLength() int {
	return request.tokenLength
}

// SetTokenType sets the token type
func (request *WidgetSessionRequest) SetTokenType(tokenType VerificationTokenType) {
	request.tokenType = tokenType
}

// GetTokenType gets the token type
func (request WidgetSessionRequest) GetTokenType() VerificationTokenType {
	return request.tokenType
}

// GetBackupCodeIdentifier gets the backup code identifier
func (request WidgetSessionRequest) GetBackupCodeIdentifier() string {
	return request.backupCodeIdentifier
}

// SetBackupCodeIdentifier sets the backup code identifier
func (request *WidgetSessionRequest) SetBackupCodeIdentifier(backupCodeIdentifier string) {
	request.backupCodeIdentifier = backupCodeIdentifier
}

// GetTotpIdentifier gets the totpCreate identifier
func (request WidgetSessionRequest) GetTotpIdentifier() string {
	return request.totpIdentifier
}

// SetTotpIdentifier sets the backup code identifier
func (request *WidgetSessionRequest) SetTotpIdentifier(totpIdentifier string) {
	request.totpIdentifier = totpIdentifier
}

// GetAllowedTypes get the allowed types
func (request *WidgetSessionRequest) GetAllowedTypes() VerificationTypes {
	return request.allowedTypes
}

// SetAllowedTypes is a shorthand to set the allowedTypes
func (request *WidgetSessionRequest) SetAllowedTypes(verifications interface{}) {
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
func (request WidgetSessionRequest) GetIssuer() string {
	return request.issuer
}

// SetIssuer sets the issuer string for line, push and telegram
func (request *WidgetSessionRequest) SetIssuer(issuer string) {
	request.issuer = issuer
}

// Submit wil submit the verification request
func (request *WidgetSessionRequest) Submit() (*WidgetSessionResponse, error) {
	response := &WidgetSessionResponse{}
	response.allowedTypes = &VerificationTypes{}
	response.requestedTypes = &VerificationTypes{}

	apiURL, err := GetURLFor("widget/session")
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
func (request WidgetSessionRequest) String() string {
	ret, _ := json.Marshal(request)
	return string(ret)
}

// WidgetSessionResponse struct that the server returns for a verification request
type WidgetSessionResponse struct {
	sessionToken           string
	applicationTag         string
	bodyTemplate           string
	createdDateTime        time.Time
	dcs                    int
	issuer                 string
	language               string
	recipient              Recipient
	sender                 string
	senderNpi              int
	senderTon              int
	tag                    string
	tokenLength            int
	tokenType              VerificationTokenType
	requestedTypes         *VerificationTypes
	allowedTypes           *VerificationTypes
	validity               int
	statusMsg              string
	statusCode             VerificationStatusCode
	salesPrice             *float32
	salesPriceCurrencyCode *string
	backupCodeIdentifier   string
	verificationIds        []string
	verification           *VerificationResponse
	links                  HATEOASLinks
}

type jsonWidgetSessionResponse struct {
	SessionToken           string                 `json:"sessionToken"`
	ApplicationTag         string                 `json:"applicationTag"`
	BodyTemplate           string                 `json:"bodyTemplate,omitempty"`
	CreatedDateTime        time.Time              `json:"createdDateTime"`
	Dcs                    int                    `json:"dcs,omitempty"`
	Issuer                 string                 `json:"issuer,omitempty"`
	Language               string                 `json:"language,omitempty"`
	Recipient              Recipient              `json:"recipient"`
	Sender                 string                 `json:"sender,omitempty"`
	SenderNpi              int                    `json:"senderNpi,omitempty"`
	SenderTon              int                    `json:"senderTon,omitempty"`
	Tag                    string                 `json:"tag,omitempty"`
	TokenLength            int                    `json:"tokenLenght,omitempty"`
	TokenType              VerificationTokenType  `json:"tokenType,omitempty"`
	RequestedTypes         *VerificationTypes     `json:"requestedTypes,omitempty"`
	AllowedTypes           *VerificationTypes     `json:"allowedTypes"`
	Validity               int                    `json:"validity"`
	StatusMsg              string                 `json:"status"`
	StatusCode             VerificationStatusCode `json:"statusCode"`
	SalesPrice             *float32               `json:"salesPrice,omitempty"`
	SalesPriceCurrencyCode *string                `json:"salesPriceCurrencyCode,omitempty"`
	BackupCodeIdentifier   string                 `json:"backupCodeIdentifier,omitempty"`
	TotpIdentifier         string                 `json:"totpIdentifier,omitempty"`
	VerificationIds        []string               `json:"verificationIds"`
	Verification           *VerificationResponse  `json:"verification,omitempty"`
	Links                  HATEOASLinks           `json:"_links"`
}

// UnmarshalJSON unmarshals the json returned by the server into the VerificationResponse struct
func (response *WidgetSessionResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonWidgetSessionResponse{}
	jsonResponse.AllowedTypes = &VerificationTypes{}
	jsonResponse.RequestedTypes = &VerificationTypes{}
	jsonResponse.Verification = &VerificationResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *WidgetSessionResponse) copyFrom(j *jsonWidgetSessionResponse) error {
	var err error // default err is nil

	response.sessionToken = j.SessionToken
	response.applicationTag = j.ApplicationTag
	response.bodyTemplate = j.BodyTemplate
	response.createdDateTime = j.CreatedDateTime
	response.dcs = j.Dcs
	response.issuer = j.Issuer
	response.language = j.Language
	response.recipient = j.Recipient
	response.sender = j.Sender
	response.senderNpi = j.SenderNpi
	response.senderTon = j.SenderTon
	response.tag = j.Tag
	response.tokenLength = j.TokenLength
	response.tokenType = j.TokenType
	response.validity = j.Validity
	response.statusMsg = j.StatusMsg
	response.statusCode = j.StatusCode
	response.salesPrice = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.verificationIds = j.VerificationIds
	response.links = j.Links
	return err
}

// GetSessionToken gets the session token
func (response WidgetSessionResponse) GetSessionToken() string {
	return response.sessionToken
}

// GetApplicationTag gets the application tag set on the verification request
func (response WidgetSessionResponse) GetApplicationTag() string {
	return response.applicationTag
}

// GetBodyTemplate gets the body templates that was set on the verification request
func (response WidgetSessionResponse) GetBodyTemplate() string {
	return response.bodyTemplate
}

// GetCreateDateTime gets the time.Time that the request was created
func (response WidgetSessionResponse) GetCreateDateTime() time.Time {
	return response.createdDateTime
}

// GetDcs gets the dcs of the verification request
func (response WidgetSessionResponse) GetDcs() int {
	return response.dcs
}

// GetLanguage gets the language that the verification request was made in (not implemented yet)
func (response WidgetSessionResponse) GetLanguage() string {
	return response.language
}

// GetRecipient get the recipient of the verification
func (response WidgetSessionResponse) GetRecipient() Recipient {
	return response.recipient
}

// GetSalesPrice get the sales price of the verification
func (response WidgetSessionResponse) GetSalesPrice() *float32 {
	return response.salesPrice
}

// GetSalesPriceCurrencyCode get the currency code of the sales price
func (response WidgetSessionResponse) GetSalesPriceCurrencyCode() *string {
	return response.salesPriceCurrencyCode
}

// GetSender get the sender
func (response WidgetSessionResponse) GetSender() string {
	return response.sender
}

// GetSenderNpi get the senderNpi
func (response WidgetSessionResponse) GetSenderNpi() int {
	return response.senderNpi
}

// GetSenderTon get the senderTon
func (response WidgetSessionResponse) GetSenderTon() int {
	return response.senderTon
}

// GetStatusMsg get the status message
func (response WidgetSessionResponse) GetStatusMsg() string {
	return response.statusMsg
}

// GetStatusCode get the status code
func (response WidgetSessionResponse) GetStatusCode() VerificationStatusCode {
	return response.statusCode
}

// GetTag get the application tag
func (response WidgetSessionResponse) GetTag() string {
	return response.tag
}

// GetTokenLength get the token length
func (response WidgetSessionResponse) GetTokenLength() int {
	return response.tokenLength
}

// GetTokenType get the token type
func (response WidgetSessionResponse) GetTokenType() VerificationTokenType {
	return response.tokenType
}

// GetValidity get the validity in seconds
func (response WidgetSessionResponse) GetValidity() int {
	return response.validity
}

// GetBackupCodeIdentifier gets the backupCodeIdentifier
func (response WidgetSessionResponse) GetBackupCodeIdentifier() string {
	return response.backupCodeIdentifier
}

// Status Retrieve the status of the validation
func (response *WidgetSessionResponse) Status() error {
	newResponse := &WidgetSessionResponse{}

	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		&response.links.Self.Href,
		nil,
		http.StatusOK,
		newResponse,
	)

	if err == nil {
		// no error use response to override ourselves
		*response = *newResponse
	}

	return err
}

// Verify verifies the token entered for a verification
func (response *WidgetSessionResponse) Verify() error {
	newResponse := &WidgetSessionResponse{}

	// to validate we need to add a query token=<token>
	newResponse.links = response.links.getDeepClone()

	q := newResponse.links.Self.Href.Query()
	q.Add("recipient", string(response.GetRecipient()))
	if response.GetBackupCodeIdentifier() != "" {
		q.Add("backupCodeIdentifier", response.GetBackupCodeIdentifier())
	}
	newResponse.links.Self.Href.RawQuery = q.Encode()

	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		&newResponse.links.Self.Href,
		nil,
		http.StatusOK,
		newResponse,
	)

	switch apiError := err.(type) {
	case *APIError:
		if apiError.Status() == http.StatusUnprocessableEntity &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenInvalid {
			response.statusCode = VerificationTokenInvalid
		} else if apiError.Status() == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenExpired {
			response.statusCode = VerificationTokenExpired
		} else if apiError.Status() == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenAlreadyVerified {
			response.statusCode = VerificationTokenAlreadyVerified
		} else if apiError.Status() == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode()) == VerificationTokenFailed {
			response.statusCode = VerificationTokenFailed
		} else {
			// undocumented response, error out
			return err
		}
		return nil
	}

	if err == nil {
		// no error use response to override ourselves
		*response = *newResponse
	}

	return err
}

// IsTokenUnknown is a helper function to check the token status is unknown
func (response WidgetSessionResponse) IsTokenUnknown() bool {
	return response.statusCode == VerificationTokenUnknown
}

// IsTokenSuccess is a helper function to check if the token status is successful
func (response WidgetSessionResponse) IsTokenSuccess() bool {
	return response.statusCode == VerificationTokenSuccess
}

// IsTokenInvalid is a helper function to check if the token status is invalid
func (response WidgetSessionResponse) IsTokenInvalid() bool {
	return response.statusCode == VerificationTokenInvalid
}

// IsTokenExpired is a helper function to check if the token has expired
func (response WidgetSessionResponse) IsTokenExpired() bool {
	return response.statusCode == VerificationTokenExpired
}

// IsTokenAlreadyVerified is a helper function to check if that token has already been verified
func (response WidgetSessionResponse) IsTokenAlreadyVerified() bool {
	return response.statusCode == VerificationTokenAlreadyVerified
}

// IsTokenFailed is a helper function to check if the token has failed to validate
func (response WidgetSessionResponse) IsTokenFailed() bool {
	return response.statusCode == VerificationTokenFailed
}

// String stringify a validation response
func (response WidgetSessionResponse) String() string {
	ret, _ := json.Marshal(response)
	return string(ret)
}

// NewWidgetSessionRequest creates a new widgetsession using a recipient (the only required var)
func NewWidgetSessionRequest() *WidgetSessionRequest {
	widgetSessionRequest := &WidgetSessionRequest{}
	widgetSessionRequest.allowedTypes = VerificationTypes{}
	return widgetSessionRequest
}

// WidgetSessionStatus retrieves the status of the validation using the messageId
func WidgetSessionStatus(sessionToken string) (*WidgetSessionResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("widget/session/%s", url.PathEscape(sessionToken)))
	if err != nil {
		return nil, err
	}

	request := &WidgetSessionResponse{sessionToken: sessionToken, links: createSelfLinks(apiURL)}
	err = request.Status()
	return request, err
}
