package twizo

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// NumberLookupRequest struct
type NumberLookupRequest struct {
	numbers     []Recipient
	tag         string
	validity    int
	resultType  ResultType
	callbackURL *url.URL // only relevant for SmsResultTypeCallback | SmsResultTYpeCallbackPollling
}

type jsonNumberLookupRequest struct {
	Numbers     []Recipient `json:"numbers"`
	Tag         string      `json:"tag,omitempty"`
	Validity    int         `json:"validity,omitempty"`
	ResultType  ResultType  `json:"resultType,omitempty"`
	CallbackURL *url.URL    `json:"callbackUrl,omitempty"`
}

// MarshalJSON is used to convert NumberLookupRequest to json
func (request *NumberLookupRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonNumberLookupRequest{
		Numbers:    request.numbers,
		Tag:        request.tag,
		Validity:   request.validity,
		ResultType: request.resultType,
	}

	// set the callback url if we need one, it still might be empty
	if request.resultType == ResultTypeCallback || request.resultType == ResultTypeCallbackPolling {
		jsonRequest.CallbackURL = request.callbackURL
	}

	return json.Marshal(jsonRequest)
}

// SetNumbers sets the numbers for a numberlookup request
func (request *NumberLookupRequest) SetNumbers(numbers []Recipient) {
	request.numbers = numbers
}

// GetNumbers gets the numbers of a numberlookup request
func (request NumberLookupRequest) GetNumbers() []Recipient {
	return request.numbers
}

// SetTag sets the tag for a numberlookup request
func (request *NumberLookupRequest) SetTag(tag string) {
	request.tag = tag
}

// GetTag returns the tag of a numberlookup request
func (request NumberLookupRequest) GetTag() string {
	return request.tag
}

// SetValidity sets the validity for a numberlookup request
func (request *NumberLookupRequest) SetValidity(validity int) {
	request.validity = validity
}

// GetValidation returns the validity of a numberlookup request
func (request NumberLookupRequest) GetValidation() int {
	return request.validity
}

// SetResultType sets the result type for a numberlookup request
func (request *NumberLookupRequest) SetResultType(resultType ResultType) {
	request.resultType = resultType
}

// GetResultType gets the requested result type of a numberlookup request
func (request NumberLookupRequest) GetResultType() ResultType {
	return request.resultType
}

// SetCallbackURL sets the callback url for a numberlookup request
func (request *NumberLookupRequest) SetCallbackURL(URL *url.URL) {
	request.callbackURL = URL
}

// GetCallbackURL gets the callback url of a numberlookup request
func (request NumberLookupRequest) GetCallbackURL() *url.URL {
	return request.callbackURL
}

// Submit actually submits the numberlookup request
func (request *NumberLookupRequest) Submit() (*NumberLookupResponses, error) {
	responses := &NumberLookupResponses{}

	apiURL, err := GetURLFor("numberlookup/submit")
	if err != nil {
		return nil, err
	}

	// todo: we need to clear our dcs and udh here, as they are not valid for simple submit
	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodPost,
		apiURL,
		request,
		http.StatusCreated,
		responses,
	)
	if err != nil {
		return nil, err
	}

	return responses, nil
}
