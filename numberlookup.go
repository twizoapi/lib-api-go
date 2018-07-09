package twizo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// NumberLookupStatusCode is actually an int
type NumberLookupStatusCode int

// Mapping of numberlookup statis codes
const (
	NumberLookupStatusCodeNoStatus    NumberLookupStatusCode = 0
	NumberLookupStatusCodeDelivered   NumberLookupStatusCode = 1
	NumberLookupStatusCodeRejected    NumberLookupStatusCode = 2
	NumberLookupStatusCodeExpired     NumberLookupStatusCode = 3
	NumberLookupStatusCodeEnroute     NumberLookupStatusCode = 4
	NumberLookupStatusCodeBuffered    NumberLookupStatusCode = 5
	NumberLookupStatusCodeAccepted    NumberLookupStatusCode = 6
	NumberLookupStatusCodeUndelivered NumberLookupStatusCode = 7
	NumberLookupStatusCodeDeleted     NumberLookupStatusCode = 8
	NumberLookupStatusCodeUnknown     NumberLookupStatusCode = 9
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

func (request *jsonNumberLookupRequest) copyFrom(r *NumberLookupRequest) {
	request.Numbers = r.numbers
	request.Tag = r.tag
	request.Validity = r.validity
	request.ResultType = r.resultType

	// set the callback url if we need one, it still might be empty
	if r.resultType == ResultTypeCallback || r.resultType == ResultTypeCallbackPolling {
		request.CallbackURL = r.callbackURL
	}
}

// MarshalJSON is used to convert NumberLookupRequest to json
func (request *NumberLookupRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := &jsonNumberLookupRequest{}
	jsonRequest.copyFrom(request)
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

// NumberLookupResponse struct
type NumberLookupResponse struct {
	applicationTag         string
	callbackURL            *url.URL
	createDateTime         time.Time
	imsi                   *string
	ported                 string
	roaming                string
	messageID              string
	msc                    *string
	networkCode            *int
	number                 string
	operator               *string
	reasonCode             *int
	resultTimestamp        string
	resultType             int
	salesPrice             *float32
	salesPriceCurrencyCode *string
	statusMsg              string
	statusCode             NumberLookupStatusCode
	tag                    *string
	validity               int
	validUntilDateTime     time.Time
	links                  HATEOASLinks
}

type jsonNumberLookupResponse struct {
	ApplicationTag         string                 `json:"applicationTag,omitempty"`
	CallbackURL            *string                `json:"callbackUrl,omitempty"`
	CountryCode            *string                `json:"countryCode,omitempty"`
	CreateDateTime         time.Time              `json:"createdDateTime,omitempty"`
	Imsi                   *string                `json:"imsi,omitempty"`
	Ported                 string                 `json:"isPorted,omitempty"`
	Roaming                string                 `json:"isRoaming,omitempty"`
	MessageID              string                 `json:"messageId"`
	Msc                    *string                `json:"msc,omitempty"`
	NetworkCode            *int                   `json:"networkCode,omitempty"`
	Number                 string                 `json:"number"`
	Operator               *string                `json:"operator"`
	ReasonCode             *int                   `json:"reasonCode,omitempty"`
	ResultTimestamp        string                 `json:"resulttimestamp,omitempty"`
	ResultType             int                    `json:"resultType,omitempty"`
	SalesPrice             *float32               `json:"salesPrice,omitempty"`
	SalesPriceCurrencyCode *string                `json:"salesPriceCurrencyCode,omitempty"`
	StatusMsg              string                 `json:"status,omitempty"`
	StatusCode             NumberLookupStatusCode `json:"statusCode,omitempty"`
	Tag                    *string                `json:"tag,omitempty"`
	Validity               int                    `json:"validity,omitempty"`
	ValidUntilDateTime     time.Time              `json:"validUntilDateTime,omitempty"`
	Links                  HATEOASLinks           `json:"_links"`
}

// UnmarshalJSON the json response to struct
func (response *NumberLookupResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonNumberLookupResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *NumberLookupResponse) copyFrom(j *jsonNumberLookupResponse) error {

	response.applicationTag = j.ApplicationTag
	response.createDateTime = j.CreateDateTime
	response.imsi = j.Imsi
	response.ported = j.Ported
	response.roaming = j.Roaming
	response.messageID = j.MessageID
	response.msc = j.Msc
	response.networkCode = j.NetworkCode
	response.number = j.Number
	response.operator = j.Operator
	response.reasonCode = j.ReasonCode
	response.resultTimestamp = j.ResultTimestamp
	response.resultType = j.ResultType
	response.salesPrice = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.statusMsg = j.StatusMsg
	response.statusCode = j.StatusCode
	response.tag = j.Tag
	response.validity = j.Validity
	response.validUntilDateTime = j.ValidUntilDateTime
	response.links = j.Links

	if j.CallbackURL != nil {
		u, err := url.Parse(*j.CallbackURL)
		if err != nil {
			return err
		}
		response.callbackURL = u
	}

	return nil
}

// Status requests the status of a numberlookup
func (response *NumberLookupResponse) Status() error {
	newNumberLookupResponse := &NumberLookupResponse{}

	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		&response.links.Self.Href,
		nil,
		http.StatusOK,
		response,
	)

	if err == nil {
		// no error use response to override ourselves
		*newNumberLookupResponse = *response
	}

	return err
}

// GetApplicationTag gets the application tag of the response
func (response NumberLookupResponse) GetApplicationTag() string {
	return response.applicationTag
}

// GetCallbackURL gets the callback url of the response
func (response NumberLookupResponse) GetCallbackURL() *url.URL {
	return response.callbackURL
}

// GetCreateDateTime gets the created time
func (response NumberLookupResponse) GetCreateDateTime() time.Time {
	return response.createDateTime
}

// GetImsi gets the imsi of the response
func (response NumberLookupResponse) GetImsi() *string {
	return response.imsi
}

// GetPorted gets ported of the response
func (response NumberLookupResponse) GetPorted() string {
	return response.ported
}

// GetRoaming gets roaming of the response
func (response NumberLookupResponse) GetRoaming() string {
	return response.roaming
}

// GetMessageID gets the message id of the response
func (response NumberLookupResponse) GetMessageID() string {
	return response.messageID
}

// GetMsc gets the msc of the response
func (response NumberLookupResponse) GetMsc() *string {
	return response.msc
}

// GetNetworkCode gets the networkCode of the response
func (response NumberLookupResponse) GetNetworkCode() *int {
	return response.networkCode
}

// GetNumber gets the number of the response
func (response NumberLookupResponse) GetNumber() string {
	return response.number
}

// GetOperator gets the operator of the response
func (response NumberLookupResponse) GetOperator() *string {
	return response.operator
}

// GetReasonCode gets the reasonCode of the response
func (response NumberLookupResponse) GetReasonCode() *int {
	return response.reasonCode
}

// GetResultType get the resultType of the response
func (response NumberLookupResponse) GetResultType() int {
	return response.resultType
}

// GetSalesPrice gets the salesprice of the response
func (response NumberLookupResponse) GetSalesPrice() *float32 {
	return response.salesPrice
}

// GetSalesPriceCurrencyCode gets the currency of the salesprice of the response
func (response NumberLookupResponse) GetSalesPriceCurrencyCode() *string {
	return response.salesPriceCurrencyCode
}

// GetStatusMsg gets the status of the response
func (response NumberLookupResponse) GetStatusMsg() string {
	return response.statusMsg
}

// GetStatusCode gets the statuscode of the response
func (response NumberLookupResponse) GetStatusCode() NumberLookupStatusCode {
	return response.statusCode
}

// GetTag gets the tag of the response
func (response NumberLookupResponse) GetTag() *string {
	return response.tag
}

// GetResultTimeStamp of the response
func (response NumberLookupResponse) GetResultTimeStamp() string {
	return response.resultTimestamp
}

// GetValidity gets the validity of the response
func (response NumberLookupResponse) GetValidity() int {
	return response.validity
}

// GetValidUntillDateTime gets the valid untill time of the response
func (response NumberLookupResponse) GetValidUntillDateTime() time.Time {
	return response.validUntilDateTime
}

type jsonNumberLookupResponseEmbedded struct {
	Items *[]NumberLookupResponse `json:"items"`
}

// Status requests the status of all embedded numberlookups
func (v *jsonNumberLookupResponseEmbedded) Status() error {
	// we potentially double memory here
	newItems := []NumberLookupResponse{}
	for _, item := range *v.Items {
		err := item.Status()
		newItems = append(newItems, item)
		if err != nil {
			return err
		}
	}
	// replace Items freeing extra memory
	v.Items = &newItems
	return nil
}

// NumberLookupResponses struct
type NumberLookupResponses struct {
	Responses *[]NumberLookupResponse
}

type jsonNumberLookupResponses struct {
	Links    HATEOASLinks                     `json:"_links"`
	Embedded jsonNumberLookupResponseEmbedded `json:"_embedded"`
	Total    int                              `json:"total_items"`
}

// UnmarshalJSON unmarshals from json to struct
func (responses *NumberLookupResponses) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonNumberLookupResponses{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	responses.Responses = jsonResponse.Embedded.Items

	return nil
}

// GetItems returns all numberlookup responses
func (responses NumberLookupResponses) GetItems() []NumberLookupResponse {
	return *responses.Responses
}

// Status requests the status of all numberlookup responses
func (responses *NumberLookupResponses) Status() error {
	newItems := []NumberLookupResponse{}
	for _, item := range *responses.Responses {
		err := item.Status()
		if err != nil {
			return err
		}
		newItems = append(newItems, item)
	}
	// all went well update Items
	responses.Responses = &newItems
	return nil
}

// NewNumberLookupRequest creates a new verificationParam using a recipient (the only required var)
func NewNumberLookupRequest(numbers []Recipient) *NumberLookupRequest {
	params := &NumberLookupRequest{
		numbers: numbers,
	}
	return params
}

// NumberLookupSubmit creates a new numberlookup and submits it
func NumberLookupSubmit(numbers interface{}) (*NumberLookupResponses, error) {
	r, err := convertRecipients(numbers)
	if err != nil {
		return nil, err
	}
	return NewNumberLookupRequest(r).Submit()
}

// NumberLookupStatus creates a new numberlookup with id and requests the status
func NumberLookupStatus(messageID string) (*NumberLookupResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("numberlookup/submit/%s", url.PathEscape(messageID)))
	if err != nil {
		return nil, err
	}

	numberLookupResponse := &NumberLookupResponse{
		messageID: messageID,
		links:     createSelfLinks(apiURL),
	}

	err = numberLookupResponse.Status()
	// were we able to find it ?
	if apiError, ok := err.(APIError); ok {
		if apiError.Status() == http.StatusNotFound {
			return nil, nil
		}
	}

	return numberLookupResponse, err
}

//
// Polling support
//
type numberLookupPollEmbedded struct {
	Messages *[]NumberLookupResponse `json:"messages"` // strange this was called items in smsResponseEmbedded
}

// NumberLookupPollResults results of the number lookup poll
type NumberLookupPollResults struct {
	pollResults
	Embedded *numberLookupPollEmbedded `json:"_embedded"`
}

func (p NumberLookupPollResults) getURL() (*url.URL, error) {
	return GetURLFor("numberlookup/poll")
}

// GetItems returns the embedded messages of a poll action
func (p NumberLookupPollResults) GetItems() []NumberLookupResponse {
	return *p.Embedded.Messages
}

// NewNumberLookupPoll creates a new NumberLookupPoll
func NewNumberLookupPoll() *NumberLookupPollResults {
	params := &NumberLookupPollResults{}
	return params
}

// NumberLookupPollStatus [todo: refactor]
func NumberLookupPollStatus() (*NumberLookupPollResults, error) {
	request := NewNumberLookupPoll()
	err := request.Status()

	return request, err
}

// Status requests the status of a NumberLookupPollResults
func (p *NumberLookupPollResults) Status() error {
	apiURL, err := p.getURL()
	if err != nil {
		return err
	}
	resp := NewNumberLookupPoll()
	err = p.status(apiURL, resp)
	if err == nil {
		*p = *resp
	}
	return err
}
