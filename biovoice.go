package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type jsonBioVoiceRequest struct {
	Recipient Recipient `json:"recipient"`
}

// BioVoiceRequest request for creating backup codes for id
type BioVoiceRequest struct {
	recipient Recipient
}

// MarshalJSON is used to convert BioVoiceRequest to json
func (request *BioVoiceRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonBioVoiceRequest{
		Recipient: request.recipient,
	}

	return json.Marshal(jsonRequest)
}

// DeleteSubscription the existing (if any) biovoice subscription
func (request *BioVoiceRequest) DeleteSubscription() error {
	apiURL, err := GetURLFor(
		fmt.Sprintf("biovoice/subscription/%s",
			url.PathEscape(string(request.GetRecipient()))),
	)
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

// CreateRegistration will trigger the creation of a boivoice registration
func (request *BioVoiceRequest) CreateRegistration() (*BioVoiceResponse, error) {
	response := &BioVoiceResponse{}

	apiURL, err := GetURLFor("biovoice/registration")
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

// CheckSubscription checks the status of the subscription
func (request *BioVoiceRequest) CheckSubscription() (*BioVoiceResponse, error) {
	response := &BioVoiceResponse{}

	return response, nil
}

// CheckRegistration checks the status of the registration
func (request *BioVoiceRequest) CheckRegistration() (*BioVoiceResponse, error) {
	response := &BioVoiceResponse{}

	return response, nil
}

// GetRecipient returns the recipient of the biovoice request
func (request *BioVoiceRequest) GetRecipient() Recipient {
	return request.recipient
}

/*
{
		"createdDateTime": "2018-02-02T15:04:15+00:00",
		"language": null,
		"reasonCode": null,
		"recipient": "6100000000",
		"registrationId": "{{ .RegistrationId }}",
		"salesPrice": null,
		"salesPriceCurrencyCode": null,
		"status": "no status",
		"statusCode": 0,
		"voiceSentence": "Verify me with my voicepin",
		"webHook": null,
		"_links": {
			"self": {
				"href": "https://{{ .Host }}/v1/biovoice/registration/{{ .RegistrationId }}"
			}
		}
	}
*/
type jsonBioVoiceResponse struct {
	CreatedDateTime        *time.Time   `json:"createdDateTime"`
	Language               *string      `json:"language"`
	ReasonCode             *string      `json:"reasonCode"`
	Recipient              Recipient    `json:"recipient"`
	RegistrationID         string       `json:"registrationId"`
	SalesPrice             *float64     `json:"salesPrice"`
	SalesPriceCurrencyCode *string      `json:"salesPriceCurrencyCode"`
	Status                 string       `json:"status"`
	StatusCode             int          `json:"statusCode"`
	VoiceSentence          string       `json:"voiceSentence"`
	WebHook                *string      `json:"webHook"`
	Links                  HATEOASLinks `json:"_links"`
}

// BioVoiceResponse response for backup create api call
type BioVoiceResponse struct {
	createdDateTime        *time.Time
	language               *string
	reasonCode             *string
	recipient              Recipient
	registrationID         string
	salesPrice             *float64
	salesPriceCurrencyCode *string
	status                 *string
	statusCode             int
	voiceSentence          string
	webHook                *string
	links                  HATEOASLinks
}

// UnmarshalJSON the json response to struct
func (response *BioVoiceResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonBioVoiceResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *BioVoiceResponse) copyFrom(j *jsonBioVoiceResponse) error {
	response.createdDateTime = j.CreatedDateTime
	response.language = j.Language
	response.reasonCode = j.ReasonCode
	response.recipient = j.Recipient
	response.registrationID = j.RegistrationID
	response.salesPrice = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.voiceSentence = j.VoiceSentence
	response.webHook = j.WebHook
	response.links = j.Links

	return nil
}

// GetRecipient gets the recipient sent with the request
func (response BioVoiceResponse) GetRecipient() Recipient {
	return response.recipient
}

//// GetCodes returns the codes that were assined (only for a create operation)
//func (response BioVoiceResponse) GetCodes() []string {
//	return response.codes
//}
//
//// GetAmountOfCodesLeft returns the amount of codes left to try
//func (response BioVoiceResponse) GetAmountOfCodesLeft() int {
//	return response.amountOfCodesLeft
//}
//
//// GetCreateDateTime returns the date timestamp when the codes were generated
//func (response BioVoiceResponse) GetCreateDateTime() *time.Time {
//	return response.createDateTime
//}
//
//// GetVerificationResponse returns the verification response or nil
//func (response *BioVoiceResponse) GetVerificationResponse() *VerificationResponse {
//	return response.verificationResponse
//}
//
//// AlreadyExists returns false if they already exist
//func (response BioVoiceResponse) AlreadyExists() bool {
//	return response.alreadyExists
//}

// NewBioVoiceRequest creates a new BioVoiceRequest
func NewBioVoiceRequest(recipient interface{}) (*BioVoiceRequest, error) {
	r, err := convertRecipients(recipient)
	if err != nil {
		return nil, err
	}
	if len(r) != 1 {
		return nil, fmt.Errorf("need at least 1 recipient got [%d]", len(r))
	}

	request := &BioVoiceRequest{recipient: r[0]}
	return request, nil
}

// BioVoiceCreateRegistration creates new biovoice registration for a recipient
func BioVoiceCreateRegistration(recipient interface{}) (*BioVoiceResponse, error) {
	request, err := NewBioVoiceRequest(recipient)
	if err != nil {
		return nil, err
	}
	return request.CreateRegistration()
}

// BioVoiceCheckRegistration checks the biovoice registration of a recipient
func BioVoiceCheckRegistration(recipient interface{}) (*BioVoiceResponse, error) {
	request, err := NewBioVoiceRequest(recipient)
	if err != nil {
		return nil, err
	}
	return request.CheckRegistration()
}

// BioVoiceCheckSubscription checks the biovoice subscription of a recipient
func BioVoiceCheckSubscription(recipient interface{}) (*BioVoiceResponse, error) {
	request, err := NewBioVoiceRequest(recipient)
	if err != nil {
		return nil, err
	}
	return request.CheckSubscription()
}

// BioVoiceDeleteSubscription will delete the biovoice for the identifier supplied
func BioVoiceDeleteSubscription(recipient interface{}) error {
	request, err := NewBioVoiceRequest(recipient)
	if err != nil {
		return err
	}
	return request.DeleteSubscription()
}
