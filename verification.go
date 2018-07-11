package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// VerificationType describes the verification type
type VerificationType string

// All verification types supported
const (
	VerificationTypeSms               VerificationType = "sms"
	VerificationTypeCall              VerificationType = "call"
	VerificationTypeBioVoice          VerificationType = "biovoice"
	VerificationTypePush              VerificationType = "push"
	VerificationTypeTotp              VerificationType = "totp"
	VerificationTypeTelegram          VerificationType = "telegram"
	VerificationTypeLine              VerificationType = "line"
	VerificationTypeBackupCode        VerificationType = "backupcode"
	VerificationTypeFacebookMessenger VerificationType = "facebook_messenger"
)

// VerificationTypes contains array of VerificationType elements
type VerificationTypes []VerificationType

// UnmarshalJSON unmarshals the json returned by the server into the VerificationTypes slice
func (vT *VerificationTypes) UnmarshalJSON(j []byte) error {
	var jsonResponse interface{}
	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	// clear out existing values
	// *vT = nil
	switch v := jsonResponse.(type) {
	case []interface{}:
		for _, element := range v {
			if e, ok := element.(string); ok {
				vT.Add(VerificationType(e))
			} else {
				return fmt.Errorf("unexpected type for VerificationType [%v]", element)
			}
		}
	case map[string]interface{}:
		// there is a bug in the api where the array is returned as a map instead of
		// an array
		for _, element := range v {
			if e, ok := element.(string); ok {
				vT.Add(VerificationType(e))
			} else {
				return fmt.Errorf("unexpected type for VerificationType [%v]", element)
			}
		}
	default:
		return fmt.Errorf("unexpected type for VerificationTypes [%#v]", v)
	}

	return nil
}

// Fetch retrieves the valid verification types for the application
func (vT *VerificationTypes) Fetch() error {
	// todo: this should be array according to documentation
	response := &VerificationTypes{}

	apiURL, _ := GetURLFor("application/verification_types")

	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		apiURL,
		nil,
		http.StatusOK,
		response,
	)

	if err == nil {
		*vT = *response
	}

	return err
}

// Add will add a verification type to
func (vT *VerificationTypes) Add(verificationType VerificationType) {
	if vT.Has(verificationType) {
		return
	}
	*vT = append(*vT, verificationType)
}

// Has will return true if verificationType is present
func (vT VerificationTypes) Has(verificationTypes ...VerificationType) bool {
	for _, b := range vT {
		for _, t := range verificationTypes {
			if b == t {
				return true
			}
		}
	}
	return false
}

// NewVerificationTypes creates a new verificationtypes
func NewVerificationTypes() *VerificationTypes {
	return &VerificationTypes{}
}

// VerificationStatusCode contains the verification status code
type VerificationStatusCode int

// VerificationTokenErrorCode describes the token error code
type VerificationTokenErrorCode int

// VerificationTokenType describes the verification token type
type VerificationTokenType string

// All verification Token types
const (
	VerificationTokenTypeDefault VerificationTokenType = ""
	VerificationTokenTypeNumeric VerificationTokenType = "numeric"
	VerificationTokenTypeAlpha   VerificationTokenType = "alphanumeric"
)

// Result of verification tokens
const (
	VerificationTokenUnknown         VerificationStatusCode = 0
	VerificationTokenSuccess         VerificationStatusCode = 1
	VerificationTokenAlreadyVerified VerificationStatusCode = 101
	VerificationTokenExpired         VerificationStatusCode = 102
	VerificationTokenInvalid         VerificationStatusCode = 103
	VerificationTokenFailed          VerificationStatusCode = 104
)

// VerificationRequest contains the struct of the request send for a verification
type VerificationRequest struct {
	recipient        Recipient
	bodyTemplate     string
	dcs              int
	language         string
	sender           string
	senderNpi        int
	senderTon        int
	sessionID        string
	tag              string
	tokenLength      int
	tokenType        VerificationTokenType
	verificationType VerificationType
	validity         string
}

type jsonVerificationRequest struct {
	Recipient        Recipient             `json:"recipient"`
	BodyTemplate     string                `json:"bodyTemplate,omitempty"`
	Dcs              int                   `json:"dcs,omitempty"`
	Language         string                `json:"language,omitempty"`
	Sender           string                `json:"sender,omitempty"`
	SenderNpi        int                   `json:"senderNpi,omitempty"`
	SenderTon        int                   `json:"senderTon,omitempty"`
	SessionID        string                `json:"sessionId,omitempty"`
	Tag              string                `json:"tag,omitempty"`
	TokenLength      int                   `json:"tokenLength,omitempty"`
	TokenType        VerificationTokenType `json:"tokenType,omitempty"`
	VerificationType VerificationType      `json:"type,omitempty"`
	Validity         string                `json:"validity,omitempty"`
}

func (request *jsonVerificationRequest) copyFrom(r *VerificationRequest) {
	request.Recipient = r.recipient
	request.BodyTemplate = r.bodyTemplate
	request.Dcs = r.dcs
	request.Language = r.language
	request.Sender = r.sender
	request.SenderTon = r.senderTon
	request.SenderNpi = r.senderNpi
	request.Tag = r.tag
	request.TokenType = r.tokenType
	request.TokenLength = r.tokenLength
	request.VerificationType = r.verificationType
	request.Validity = r.validity
}

// MarshalJSON is used to convert SmsRequest to json
func (request *VerificationRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := &jsonVerificationRequest{}
	jsonRequest.copyFrom(request)
	return json.Marshal(jsonRequest)
}

// SetRecipient set the recipient of the verification
func (request *VerificationRequest) SetRecipient(recipient Recipient) {
	request.recipient = recipient
}

// GetRecipient get the recipient of the verification
func (request VerificationRequest) GetRecipient() Recipient {
	return request.recipient
}

// SetBodyTemplate set the body template %token% will be replaced by the actual token sent
func (request *VerificationRequest) SetBodyTemplate(template string) {
	request.bodyTemplate = template
}

// GetBodyTemplate get the body template
func (request VerificationRequest) GetBodyTemplate() string {
	return request.bodyTemplate
}

// SetDcs sets the dcs of the verification request
func (request *VerificationRequest) SetDcs(dcs int) {
	request.dcs = dcs
}

// GetDcs get the dcs of the verification request
func (request VerificationRequest) GetDcs() int {
	return request.dcs
}

// SetLanguage Api does not implement this yet
func (request *VerificationRequest) SetLanguage(lang string) {
	request.language = lang
}

// GetLanguage get the language (not implemented yet)
func (request VerificationRequest) GetLanguage() string {
	return request.language
}

// SetSender set the sender
func (request *VerificationRequest) SetSender(sender string) {
	request.sender = sender
}

// GetSender get the sender
func (request VerificationRequest) GetSender() string {
	return request.sender
}

// SetSenderNpi sets the senderNpi
func (request *VerificationRequest) SetSenderNpi(npi int) {
	request.senderNpi = npi
}

// GetSenderNpi gets the senderNpi
func (request VerificationRequest) GetSenderNpi() int {
	return request.senderNpi
}

// SetSenderTon sets the senderTon
func (request *VerificationRequest) SetSenderTon(ton int) {
	request.senderTon = ton
}

// GetSenderTon get the senderTon
func (request VerificationRequest) GetSenderTon() int {
	return request.senderTon
}

// SetTag set the tag for this verification
func (request *VerificationRequest) SetTag(tag string) {
	request.tag = tag
}

// GetTag gets the tag set for this verification
func (request VerificationRequest) GetTag() string {
	return request.tag
}

// SetTokenLength sets the token length
func (request *VerificationRequest) SetTokenLength(length int) {
	request.tokenLength = length
}

// GetTokenLength gets the token length
func (request VerificationRequest) GetTokenLength() int {
	return request.tokenLength
}

// SetTokenType sets the token type
func (request *VerificationRequest) SetTokenType(tokenType VerificationTokenType) {
	request.tokenType = tokenType
}

// GetTokenType gets the token type
func (request VerificationRequest) GetTokenType() VerificationTokenType {
	return request.tokenType
}

// SetVerificationType sets the verification type
func (request *VerificationRequest) SetVerificationType(verificationType VerificationType) {
	request.verificationType = verificationType
}

// GetVerificationType returns the verification type
func (request VerificationRequest) GetVerificationType() VerificationType {
	return request.verificationType
}

// SetValidity set the validity of the verification request, as amount of seconds
// todo: this should be integer
func (request *VerificationRequest) SetValidity(validity string) {
	request.validity = validity
}

// GetValidity gets the validity of a verification
func (request VerificationRequest) GetValidity() string {
	return request.validity
}

// Submit wil submit the verification request
func (request *VerificationRequest) Submit() (*VerificationResponse, error) {
	response := &VerificationResponse{}

	apiURL, err := GetURLFor("verification/submit")
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
func (request VerificationRequest) String() string {
	ret, _ := json.Marshal(request)
	return string(ret)
}

// VerificationResponse struct that the server returns for a verification request
type VerificationResponse struct {
	stdVerificationResponse
	applicationTag         string
	bodyTemplate           string
	createdDateTime        time.Time
	dcs                    int
	issuer                 *string
	language               string
	messageID              string
	reasonCode             string
	recipient              Recipient
	salesPrice             *float32
	salesPriceCurrencyCode *string
	sender                 string
	senderNpi              int
	senderTon              int
	sessionID              string
	statusMsg              string
	tag                    string
	tokenLength            string
	tokenType              string
	verificationType       string
	validity               int
	validUntilDateTime     time.Time
	voiceSentence          *string
	webHook                *string
	links                  HATEOASLinks
}

type jsonVerificationResponse struct {
	ApplicationTag         string                 `json:"applicationTag"`
	BodyTemplate           string                 `json:"bodyTemplate,omitempty"`
	CreatedDateTime        time.Time              `json:"createdDateTime"`
	Dcs                    int                    `json:"dcs,omitempty"`
	Issuer                 *string                `json:"issuer,omitempty"`
	Language               string                 `json:"language,omitempty"`
	MessageID              string                 `json:"messageId"`
	ReasonCode             string                 `json:"reasonCode,omitempty"`
	Recipient              Recipient              `json:"recipient"`
	SalesPrice             *float32               `json:"salesPrice,omitempty"`
	SalesPriceCurrencyCode *string                `json:"salesPriceCurrencyCode,omitempty"`
	Sender                 string                 `json:"sender,omitempty"`
	SenderNpi              int                    `json:"senderNpi,omitempty"`
	SenderTon              int                    `json:"senderTon,omitempty"`
	SessionID              string                 `json:"sessionId"`
	StatusMsg              string                 `json:"status"`
	StatusCode             VerificationStatusCode `json:"statusCode"`
	Tag                    string                 `json:"tag,omitempty"`
	TokenLength            string                 `json:"tokenLenght,omitempty"`
	TokenType              string                 `json:"tokenType,omitempty"`
	VerificationType       string                 `json:"type"`
	Validity               int                    `json:"validity,omitempty"`
	ValidUntilDateTime     time.Time              `json:"validUntilDateTime,omitempty"`
	VoiceSentence          *string                `json:"voiceSentence,omitempty"`
	WebHook                *string                `json:"webHook,omitempty"`
	Links                  HATEOASLinks           `json:"_links"`
}

// UnmarshalJSON unmarshals the json returned by the server into the VerificationResponse struct
func (response *VerificationResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonVerificationResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *VerificationResponse) copyFrom(j *jsonVerificationResponse) error {
	var err error // default err is nil

	response.applicationTag = j.ApplicationTag
	response.bodyTemplate = j.BodyTemplate
	response.createdDateTime = j.CreatedDateTime
	response.dcs = j.Dcs
	response.issuer = j.Issuer
	response.language = j.Language
	response.messageID = j.MessageID
	response.reasonCode = j.ReasonCode
	response.recipient = j.Recipient
	response.salesPrice = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.sender = j.Sender
	response.senderNpi = j.SenderNpi
	response.senderTon = j.SenderTon
	response.sessionID = j.SessionID
	response.statusMsg = j.StatusMsg
	response.statusCode = j.StatusCode
	response.tag = j.Tag
	response.tokenLength = j.TokenLength
	response.tokenType = j.TokenType
	response.verificationType = j.VerificationType
	response.validity = j.Validity
	response.validUntilDateTime = j.ValidUntilDateTime
	response.voiceSentence = j.VoiceSentence
	response.webHook = j.WebHook
	response.links = j.Links

	return err
}

// GetApplicationTag gets the application tag set on the verification request
func (response VerificationResponse) GetApplicationTag() string {
	return response.applicationTag
}

// GetBodyTemplate gets the body templates that was set on the verification request
func (response VerificationResponse) GetBodyTemplate() string {
	return response.bodyTemplate
}

// GetCreateDateTime gets the time.Time that the request was created
func (response VerificationResponse) GetCreateDateTime() time.Time {
	return response.createdDateTime
}

// GetDcs gets the dcs of the verification request
func (response VerificationResponse) GetDcs() int {
	return response.dcs
}

// GetLanguage gets the language that the verification request was made in (not implemented yet)
func (response VerificationResponse) GetLanguage() string {
	return response.language
}

// GetMessageID gets the message Id of the verification request
func (response VerificationResponse) GetMessageID() string {
	return response.messageID
}

// GetReasonCode get the reason code of the verification request
func (response VerificationResponse) GetReasonCode() string {
	return response.reasonCode
}

// GetRecipient get the recipient of the verification
func (response VerificationResponse) GetRecipient() Recipient {
	return response.recipient
}

// GetSalesPrice get the sales price of the verification
func (response VerificationResponse) GetSalesPrice() *float32 {
	return response.salesPrice
}

// GetSalesPriceCurrencyCode get the currency code of the sales price
func (response VerificationResponse) GetSalesPriceCurrencyCode() *string {
	return response.salesPriceCurrencyCode
}

// GetSender get the sender
func (response VerificationResponse) GetSender() string {
	return response.sender
}

// GetSenderNpi get the senderNpi
func (response VerificationResponse) GetSenderNpi() int {
	return response.senderNpi
}

// GetSenderTon get the senderTon
func (response VerificationResponse) GetSenderTon() int {
	return response.senderTon
}

// GetSessionID get the session id
func (response VerificationResponse) GetSessionID() string {
	return response.sessionID
}

// GetStatusMsg get the status message
func (response VerificationResponse) GetStatusMsg() string {
	return response.statusMsg
}

// GetStatusCode get the status code
func (response VerificationResponse) GetStatusCode() VerificationStatusCode {
	return response.statusCode
}

// GetTag get the application tag
func (response VerificationResponse) GetTag() string {
	return response.tag
}

// GetTokenLength get the token length
func (response VerificationResponse) GetTokenLength() string {
	return response.tokenLength
}

// GetTokenType get the token type
func (response VerificationResponse) GetTokenType() string {
	return response.tokenLength
}

// GetVerificationType get the verification type
func (response VerificationResponse) GetVerificationType() string {
	return response.verificationType
}

// GetValidity get the validity in seconds
func (response VerificationResponse) GetValidity() int {
	return response.validity
}

// GetValidUntilDateTime get the validity in time.Time
func (response VerificationResponse) GetValidUntilDateTime() time.Time {
	return response.validUntilDateTime
}

// Status Retrieve the status of the validation
func (response *VerificationResponse) Status() error {
	newResponse := &VerificationResponse{}

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
func (response *VerificationResponse) Verify(token string) error {
	newResponse := &VerificationResponse{}

	// to validate we need to add a query token=<token>
	newResponse.links = response.links.getDeepClone()

	q := newResponse.links.Self.Href.Query()
	q.Add("token", token)
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
func (response VerificationResponse) IsTokenUnknown() bool {
	return response.statusCode == VerificationTokenUnknown
}

// IsTokenSuccess is a helper function to check if the token status is successful
func (response VerificationResponse) IsTokenSuccess() bool {
	return response.statusCode == VerificationTokenSuccess
}

// IsTokenInvalid is a helper function to check if the token status is invalid
func (response VerificationResponse) IsTokenInvalid() bool {
	return response.statusCode == VerificationTokenInvalid
}

// IsTokenExpired is a helper function to check if the token has expired
func (response VerificationResponse) IsTokenExpired() bool {
	return response.statusCode == VerificationTokenExpired
}

// IsTokenAlreadyVerified is a helper function to check if that token has already been verified
func (response VerificationResponse) IsTokenAlreadyVerified() bool {
	return response.statusCode == VerificationTokenAlreadyVerified
}

// IsTokenFailed is a helper function to check if the token has failed to validate
func (response VerificationResponse) IsTokenFailed() bool {
	return response.statusCode == VerificationTokenFailed
}

// String stringify a validation response
func (response VerificationResponse) String() string {
	ret, _ := json.Marshal(response)
	return string(ret)
}

type stdVerificationResponse struct {
	statusCode VerificationStatusCode
}

func (response stdVerificationResponse) IsTokenSuccess() bool {
	return response.statusCode == VerificationTokenSuccess
}

// NewVerificationRequest creates a new verificationParam using a recipient (the only required var)
func NewVerificationRequest(recipient interface{}) (*VerificationRequest, error) {
	r, err := convertRecipients(recipient)
	if err != nil {
		return nil, err
	}
	if len(r) != 1 {
		return nil, fmt.Errorf("need exactly one [recipient] for NewVerificationRequest got [%d]", len(r))
	}
	params := &VerificationRequest{recipient: r[0]}
	return params, nil
}

// VerificationSubmit creates a verificationRequest from verificationParams and submits it
func VerificationSubmit(recipient interface{}) (*VerificationResponse, error) {
	verification, err := NewVerificationRequest(recipient)
	if err != nil {
		return nil, err
	}
	return verification.Submit()
}

// VerificationStatus retrieves the status of the validation using the messageId
func VerificationStatus(messageID string) (*VerificationResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("verification/submit/%s", url.PathEscape(messageID)))
	if err != nil {
		return nil, err
	}

	request := &VerificationResponse{messageID: messageID, links: createSelfLinks(apiURL)}
	err = request.Status()
	return request, err
}

// VerificationVerify validates the result of the validation request using the messageId and the token
func VerificationVerify(messageID string, token string) (*VerificationResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("verification/submit/%s", url.PathEscape(messageID)))
	if err != nil {
		return nil, err
	}

	request := &VerificationResponse{messageID: messageID, links: createSelfLinks(apiURL)}
	err = request.Verify(token)
	return request, err
}

// VerificationFetchTypes retrieves all verification types for the application
// from the server
func VerificationFetchTypes() (*VerificationTypes, error) {
	v := &VerificationTypes{}
	err := v.Fetch()
	if err != nil {
		return nil, err
	}
	return v, nil
}
