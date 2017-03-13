package twizo

import (
	"net/http"
	"encoding/json"
)

type VerificationRequest struct {
	recipient    	 Recipient
	bodyTemplate 	 string
	dcs          	 int
	language     	 string
	sender       	 string
	senderNpi    	 int
	senderTon    	 int
	sessionID    	 string
	tag          	 string
	tokenLength  	 int
	tokenType    	 VerificationTokenType
	verificationType VerificationType
	validity     	 string
}

type jsonVerificationRequest struct {
	Recipient    	 Recipient             `json:"recipient"`
	BodyTemplate 	 string                `json:"bodyTemplate,omitempty"`
	Dcs          	 int                   `json:"dcs,omitempty"`
	Language     	 string                `json:"language,omitempty"`
	Sender       	 string                `json:"sender,omitempty"`
	SenderNpi    	 int                   `json:"senderNpi,omitempty"`
	SenderTon    	 int                   `json:"senderTon,omitempty"`
	SessionID    	 string                `json:"sessionId,omitempty"`
	Tag          	 string                `json:"tag,omitempty"`
	TokenLength  	 int                   `json:"tokenLength,omitempty"`
	TokenType    	 VerificationTokenType `json:"tokenType,omitempty"`
	VerificationType VerificationType      `json:"type,omitempty"`
	Validity     	 string                `json:"validity,omitempty"`
}

// MarshalJSON is used to convert SmsRequest to json
func (request *VerificationRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonVerificationRequest {
		Recipient 	  : request.recipient,
		BodyTemplate      : request.bodyTemplate,
		Dcs               : request.dcs,
		Language          : request.language,
		Sender 		  : request.sender,
		SenderTon         : request.senderTon,
		SenderNpi         : request.senderNpi,
		Tag               : request.tag,
		TokenType         : request.tokenType,
		TokenLength       : request.tokenLength,
		VerificationType  : request.verificationType,
		Validity          : request.validity,
	}

	return json.Marshal(jsonRequest)
}

func (request *VerificationRequest) SetRecipient(recipient Recipient) {
	request.recipient = recipient
}

func (request VerificationRequest) GetRecipient() (Recipient) {
	return request.recipient
}

func (request *VerificationRequest) SetBodyTemplate(template string) {
	request.bodyTemplate = template
}

func (request VerificationRequest) GetBodyTemplate() (string) {
	return request.bodyTemplate
}

func (request *VerificationRequest) SetDcs(dcs int) {
	request.dcs = dcs
}

func (request VerificationRequest) GetDcs() (int) {
	return request.dcs
}

// SetLanguage Api does not implement this yet
func (request *VerificationRequest) SetLanguage(lang string) {
	request.language = lang
}

func (request VerificationRequest) GetLanguage() (string) {
	return request.language
}

func (request *VerificationRequest) SetSender(sender string) {
	request.sender = sender
}

func (request VerificationRequest) GetSender() (string) {
	return request.sender
}

func (request *VerificationRequest) SetSenderNpi(npi int) {
	request.senderNpi = npi
}

func (request VerificationRequest) GetSenderNpi() (int) {
	return request.senderNpi
}

func (request *VerificationRequest) SetSenderTon(ton int) {
	request.senderTon = ton
}

func (request VerificationRequest) GetSenderTon() (int) {
	return request.senderTon
}

func (request *VerificationRequest) SetTag(tag string) {
	request.tag = tag
}

func (request VerificationRequest) GetTag() (string) {
	return request.tag
}

func (request *VerificationRequest) SetTokenLength(length int) {
	request.tokenLength = length
}

func (request VerificationRequest) GetTokenLength() (int) {
	return request.tokenLength
}

func (request *VerificationRequest) SetTokenType(tokenType VerificationTokenType) {
	request.tokenType = tokenType
}

func (request VerificationRequest) GetTokenType() (VerificationTokenType) {
	return request.tokenType
}

func (request *VerificationRequest) SetVerificationType(verificationType VerificationType) {
	request.verificationType = verificationType
}

func (request VerificationRequest) GetVerificationType() (VerificationType) {
	return request.verificationType
}

func (request *VerificationRequest) SetValidity(validity string) {
	request.validity = validity
}

func (request VerificationRequest) GetValidity() (string) {
	return request.validity
}

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

func (request VerificationRequest) String() string {
	ret, _ := json.Marshal(request)
	return string(ret)
}
