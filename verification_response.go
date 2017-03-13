package twizo

import (
	"time"
	"net/http"
	"encoding/json"
)

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

func (r *VerificationResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonVerificationResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	r.copyFrom(jsonResponse)

	return nil
}

func (response *VerificationResponse) copyFrom (j *jsonVerificationResponse) (error) {
	var err error  // default err is nil

	response.applicationTag         = j.ApplicationTag
	response.bodyTemplate           = j.BodyTemplate
	response.createdDateTime        = j.CreatedDateTime
	response.dcs                    = j.Dcs
	response.language               = j.Language
	response.messageID              = j.MessageID
	response.reasonCode             = j.ReasonCode
	response.recipient              = j.Recipient
	response.salesPrice             = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.sender                 = j.Sender
	response.senderNpi              = j.SenderNpi
	response.senderTon              = j.SenderTon
	response.sessionID              = j.SessionID
	response.statusMsg              = j.StatusMsg
	response.statusCode             = j.StatusCode
	response.tag                    = j.Tag
	response.tokenLength            = j.TokenLength
	response.tokenType              = j.TokenType
	response.verificationType       = j.VerificationType
	response.validity               = j.Validity
	response.validUntilDateTime     = j.ValidUntilDateTime
	response.links                  = j.Links

	return err
}

func (response VerificationResponse) GetApplicationTag() (string) {
	return response.applicationTag
}

func (response VerificationResponse) GetBodyTemplate() (string) {
	return response.bodyTemplate
}

func (response VerificationResponse) GetCreateDateTime() (time.Time) {
	return response.createdDateTime
}

func (response VerificationResponse) GetDcs() (int) {
	return response.dcs
}

func (response VerificationResponse) GetLanguage() (string) {
	return response.language
}

func (response VerificationResponse) GetMessageID() (string) {
	return response.messageID
}

func (response VerificationResponse) GetReasonCode() (string) {
	return response.reasonCode
}

func (response VerificationResponse) GetRecipient() (Recipient) {
	return response.recipient
}

func (response VerificationResponse) GetSalesPrice() (*float32) {
	return response.salesPrice
}

func (response VerificationResponse) GetSalesPriceCurrencyCode() (*string) {
	return response.salesPriceCurrencyCode
}

func (response VerificationResponse) GetSender() (string) {
	return response.sender
}

func (response VerificationResponse) GetSenderNpi() (int) {
	return response.senderNpi
}

func (response VerificationResponse) GetSenderTon() (int) {
	return response.senderTon
}

func (response VerificationResponse) GetSessionID() (string) {
	return response.sessionID
}

func (response VerificationResponse) GetStatusMsg() (string) {
	return response.statusMsg
}

func (response VerificationResponse) GetStatusCode() (VerificationStatusCode) {
	return response.statusCode
}

func (response VerificationResponse) GetTag() (string) {
	return response.tag
}

func (response VerificationResponse) GetTokenLength() (string) {
	return response.tokenLength
}

func (response VerificationResponse) GetTokenType() (string) {
	return response.tokenLength
}

func (response VerificationResponse) GetVerificationType() (string) {
	return response.verificationType
}

func (response VerificationResponse) GetValidity() (int) {
	return response.validity
}

func (response VerificationResponse) GetValidUntilDateTime() (time.Time) {
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

// Verify
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
		if apiError.Status == http.StatusUnprocessableEntity && VerificationStatusCode(apiError.ErrorCode) == VerificationTokenInvalid {
			response.statusCode = VerificationTokenInvalid
		} else if apiError.Status == http.StatusLocked && VerificationStatusCode(apiError.ErrorCode) == VerificationTokenExpired {
			response.statusCode = VerificationTokenExpired
		} else if apiError.Status == http.StatusLocked && VerificationStatusCode(apiError.ErrorCode) == VerificationTokenAlreadyVerified {
			response.statusCode = VerificationTokenAlreadyVerified
		} else if apiError.Status == http.StatusLocked && VerificationStatusCode(apiError.ErrorCode) == VerificationTokenFailed {
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

func (response VerificationResponse) IsTokenUnknown() bool {
	return (response.statusCode == VerificationTokenUnknown)
}

func (response VerificationResponse) IsTokenSuccess() bool {
	return (response.statusCode == VerificationTokenSuccess)
}

func (response VerificationResponse) IsTokenInvalid() bool {
	return (response.statusCode == VerificationTokenInvalid)
}

func (response VerificationResponse) IsTokenExpired() bool {
	return (response.statusCode == VerificationTokenExpired)
}

func (response VerificationResponse) IsTokenAlreadyVerified() bool {
	return (response.statusCode == VerificationTokenAlreadyVerified)
}

func (response VerificationResponse) IsTokenFailed() bool {
	return (response.statusCode == VerificationTokenFailed)
}

/**
 * Stringify a validation response
 */
func (response VerificationResponse) String() string {
	ret, _ := json.Marshal(response)
	return string(ret)
}
