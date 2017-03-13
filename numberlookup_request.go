package twizo

import (
	"net/url"
	"net/http"
	"encoding/json"
)

type NumberLookupRequest struct {
	numbers     []Recipient
	tag         string
	validity    int
	resultType  ResultType
	callbackURL *url.URL		// only relevant for SmsResultTypeCallback | SmsResultTYpeCallbackPollling
}

type jsonNumberLookupRequest struct {
	Numbers     []Recipient            `json:"numbers"`
	Tag         string                 `json:"tag,omitempty"`
	Validity    int                    `json:"validity,omitempty"`
	ResultType  ResultType             `json:"resultType,omitempty"`
	CallbackURL *url.URL               `json:"callbackUrl,omitempty"`
}

// MarshalJSON is used to convert NumberLookupRequest to json
func (request *NumberLookupRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonNumberLookupRequest{
		Numbers           : request.numbers,
		Tag               : request.tag,
		Validity          : request.validity,
		ResultType        : request.resultType,
	}

	// set the callback url if we need one, it still might be empty
	if request.resultType == ResultTypeCallback || request.resultType == ResultTypeCallbackPolling {
		jsonRequest.CallbackURL = request.callbackURL
	}

	return json.Marshal(jsonRequest)
}

func (v *NumberLookupRequest) SetNumbers(numbers []Recipient) {
	v.numbers = numbers
}

func (v NumberLookupRequest) GetNumbers() ([]Recipient) {
	return v.numbers
}

func (v *NumberLookupRequest) SetTag(tag string) {
	v.tag = tag
}

func (v NumberLookupRequest) GetTag() (string) {
	return v.tag
}

func (v *NumberLookupRequest) SetValidity(validity int) {
	v.validity = validity
}

func (v NumberLookupRequest) GetValidation() (int) {
	return v.validity
}

func (v *NumberLookupRequest) SetResultType(resultType ResultType) {
	v.resultType = resultType
}

func (v NumberLookupRequest) GetResultType() (ResultType) {
	return v.resultType
}

func (v *NumberLookupRequest) SetCallbackUrl(URL *url.URL) {
	v.callbackURL = URL
}

func (v NumberLookupRequest) GetCallbackUrl() (*url.URL) {
	return v.callbackURL
}

func (v *NumberLookupRequest) Submit() (*NumberLookupResponses, error) {
	responses := &NumberLookupResponses{}

	apiURL, err := GetURLFor("numberlookup/submit")
	if err != nil {
		return nil, err
	}

	// todo: we need to clear our dcs and udh here, as they are not valid for simple submit
	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodPost,
		apiURL,
		v,
		http.StatusCreated,
		responses,
	)
	if err != nil {
		return nil, err
	}

	return responses, nil
}
