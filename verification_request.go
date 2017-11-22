package twizo

import (
	"encoding/json"
	"net/http"
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

// MarshalJSON is used to convert SmsRequest to json
func (request *VerificationRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonVerificationRequest{
		Recipient:        request.recipient,
		BodyTemplate:     request.bodyTemplate,
		Dcs:              request.dcs,
		Language:         request.language,
		Sender:           request.sender,
		SenderTon:        request.senderTon,
		SenderNpi:        request.senderNpi,
		Tag:              request.tag,
		TokenType:        request.tokenType,
		TokenLength:      request.tokenLength,
		VerificationType: request.verificationType,
		Validity:         request.validity,
	}

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
