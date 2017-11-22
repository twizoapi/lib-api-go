package twizo

import (
	"encoding/json"
	"net/http"
	"time"
)

// VerificationResponse struct that the server returns for a verification request
type VerificationResponse struct {
	applicationTag         string
	bodyTemplate           string
	createdDateTime        time.Time
	dcs                    int
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
	statusCode             VerificationStatusCode
	tag                    string
	tokenLength            string
	tokenType              string
	verificationType       string
	validity               int
	validUntilDateTime     time.Time
	links                  HATEOASLinks
}

type jsonVerificationResponse struct {
	ApplicationTag         string                 `json:"applicationTag"`
	BodyTemplate           string                 `json:"bodyTemplate,omitempty"`
	CreatedDateTime        time.Time              `json:"createdDateTime"`
	Dcs                    int                    `json:"dcs,omitempty"`
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
		if apiError.Status == http.StatusUnprocessableEntity &&
			VerificationStatusCode(apiError.ErrorCode) == VerificationTokenInvalid {
			response.statusCode = VerificationTokenInvalid
		} else if apiError.Status == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode) == VerificationTokenExpired {
			response.statusCode = VerificationTokenExpired
		} else if apiError.Status == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode) == VerificationTokenAlreadyVerified {
			response.statusCode = VerificationTokenAlreadyVerified
		} else if apiError.Status == http.StatusLocked &&
			VerificationStatusCode(apiError.ErrorCode) == VerificationTokenFailed {
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
