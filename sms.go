package twizo

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type smsSubmitType int

const (
	// SmsSubmitTypeSimple (default) use simple submit
	// (auto body splitting auto udh / dcs setting)
	SmsSubmitTypeSimple smsSubmitType = 0

	// SmsSubmitTypeAdvanced will not autosplit and can
	// also send binary messages
	SmsSubmitTypeAdvanced smsSubmitType = 1
)

// SmsStatusCode the status code of the sms
type SmsStatusCode int

const (
	// SmsStatusCodeNoStatus no status for message
	SmsStatusCodeNoStatus SmsStatusCode = 0

	// SmsStatusCodeDelivered message delivered
	SmsStatusCodeDelivered SmsStatusCode = 1

	// SmsStatusCodeRejected message rejected
	SmsStatusCodeRejected SmsStatusCode = 2

	// SmsStatusCodeExpired message expired
	SmsStatusCodeExpired SmsStatusCode = 3

	// SmsStatusCodeEnroute message enroute
	SmsStatusCodeEnroute SmsStatusCode = 4

	// SmsStatusCodeBuffered message buffered
	SmsStatusCodeBuffered SmsStatusCode = 5

	// SmsStatusCodeAccepted message accepted
	SmsStatusCodeAccepted SmsStatusCode = 6

	// SmsStatusCodeUndelivered message undelivered
	SmsStatusCodeUndelivered SmsStatusCode = 7

	// SmsStatusCodeDeleted message deleted
	SmsStatusCodeDeleted SmsStatusCode = 8

	// SmsStatusCodeUnknown message status unkown
	SmsStatusCodeUnknown SmsStatusCode = 9
)

// SmsRequest is used to send an sms request
type SmsRequest struct {
	recipients        []Recipient
	body              []byte        // 10 x 160
	submitType        smsSubmitType // type simple (default) or advanced (not sent to server)
	sender            string
	senderTon         int // posivive small int
	senderNpi         int // positive small int
	pid               int
	scheduledDelivery string
	tag               string
	validity          string // duration.Seconds
	resultType        ResultType
	callbackURL       *url.URL // only relevant for SmsResultTypeCallback | SmsResultTYpeCallbackPollling
	dcs               int      // relevant for advanced, 0-255
	udh               *string  // relevant for advanced, hex
}

type jsonSmsRequest struct {
	Recipients        []Recipient `json:"recipients"`
	Body              string      `json:"body"`
	Sender            string      `json:"sender"`
	SenderTon         int         `json:"senderTon,omitempty"`
	SenderNpi         int         `json:"senderNpi,omitempty"`
	Pid               int         `json:"pid,omitempty"`
	ScheduledDelivery string      `json:"scheduledDelivery,omitempty"`
	Tag               string      `json:"tag,omitempty"`
	Validity          string      `json:"validity,omitempty"`
	ResultType        ResultType  `json:"resultType"`
	CallbackURL       *url.URL    `json:"callbackUrl,omitempty"`
	Dcs               int         `json:"dcs,omitempty"`
	Udh               *string     `json:"udh,omitempty"`
}

// MarshalJSON is used to convert SmsRequest to json
func (request *SmsRequest) MarshalJSON() ([]byte, error) {
	jsonRequest := jsonSmsRequest{
		Recipients:        request.recipients,
		Body:              string(request.body),
		Sender:            request.sender,
		SenderTon:         request.senderTon,
		SenderNpi:         request.senderNpi,
		Pid:               request.pid,
		ScheduledDelivery: request.scheduledDelivery,
		Validity:          request.validity,
		ResultType:        request.resultType,
		Tag:               request.tag,
	}

	// set the callback url if we need one, it still might be empty
	if request.resultType == ResultTypeCallback || request.resultType == ResultTypeCallbackPolling {
		jsonRequest.CallbackURL = request.callbackURL
	}

	if request.submitType == SmsSubmitTypeAdvanced {
		// if request is advanced, send dcs and udh too
		jsonRequest.Dcs = request.dcs
		jsonRequest.Udh = request.udh

		// if the request is binary we hex encode the body
		if request.IsBinary() {
			jsonRequest.Body = hex.EncodeToString(request.body)
		}
	}

	return json.Marshal(jsonRequest)
}

// IsBinary returns if message is binary or not
func (request SmsRequest) IsBinary() bool {
	return isDcsBinary(request.dcs)
}

// SetSender sets the sender for an sms request
func (request *SmsRequest) SetSender(sender string) error {
	request.sender = sender

	return nil
}

// GetSender returns the sender for an sms request
func (request SmsRequest) GetSender() (string, error) {
	return request.sender, nil
}

// SetTag sets a tag for an sms request
func (request *SmsRequest) SetTag(tag string) error {
	request.tag = tag

	return nil
}

// GetTag gets the tag for an sms request
func (request SmsRequest) GetTag() (string, error) {
	return request.tag, nil
}

// SetBody set the body
func (request *SmsRequest) SetBody(body interface{}) error {
	switch t := body.(type) {
	case string:
		return request.SetBodyFromString(t)
	case []byte:
		return request.SetBodyFromByteArr(t)
	default:
		return fmt.Errorf("SetBody expects string or []byte, got [%s]", reflect.TypeOf(body))
	}
}

// SetBodyFromString set the body as string
func (request *SmsRequest) SetBodyFromString(body string) error {
	request.body = []byte(body)
	return nil
}

// SetBodyFromByteArr set the body as string
func (request *SmsRequest) SetBodyFromByteArr(body []byte) error {
	request.body = body
	return nil
}

// GetBodyAsString get the body as string
func (request SmsRequest) GetBodyAsString() (string, error) {
	return string(request.body), nil
}

// GetBodyAsByteArr get the body as []byte
func (request SmsRequest) GetBodyAsByteArr() ([]byte, error) {
	return request.body, nil
}

// SetRecipients set the recipients
func (request *SmsRequest) SetRecipients(recipients []Recipient) error {
	request.recipients = recipients

	return nil
}

// SetResultType set the resultType
func (request *SmsRequest) SetResultType(resultType ResultType) error {
	request.resultType = resultType

	return nil
}

// GetResultType get the resultType
func (request SmsRequest) GetResultType() (ResultType, error) {
	return request.resultType, nil
}

// SetSenderTon will set the senderTon of the SMS, this value will be automatically detected
// when sending the sms.  see: https://www.twizo.com/tutorials/#sender
func (request *SmsRequest) SetSenderTon(ton int) {
	request.senderTon = ton
}

// GetSenderTon retrieves the senderTon set
func (request SmsRequest) GetSenderTon() int {
	return request.senderTon
}

// SetSenderNpi set the senderNpi when sending the SMS, this value will be automatically detected
// when sending the sms. see: https://www.twizo.com/tutorials/#sender
func (request *SmsRequest) SetSenderNpi(npi int) {
	request.senderNpi = npi
}

// GetSenderNpi will get the senderNpi of the message
func (request SmsRequest) GetSenderNpi() int {
	return request.senderNpi
}

// SetPid sets an optional internal parameter allowing the sending of hidden sms (among other options)
// see: https://en.wikipedia.org/wiki/GSM_03.40#Protocol_Identifier
func (request *SmsRequest) SetPid(pid int) {
	request.pid = pid
}

// GetPid gets the Protocol Identifier of the sms (this is usually not set)
func (request SmsRequest) GetPid() int {
	return request.pid
}

// SetScheduledDelivery this will set the delivery schedule (ie when the SMS will be sent), this should must be
// in ISO-8601 format
func (request *SmsRequest) SetScheduledDelivery(sd string) {
	request.scheduledDelivery = sd
}

// GetScheduledDelivery gets the scheduled delivery
func (request SmsRequest) GetScheduledDelivery() string {
	return request.scheduledDelivery
}

// SetValidity sets the validity of the message, this value is the amount of seconds that the message is valid
// after sending the SMS, if this expires the message will be expire and no more attempts will be made
// todo: validity should be an integer
func (request *SmsRequest) SetValidity(validity string) {
	request.validity = validity
}

// GetValidity gets the validity of the message.
func (request SmsRequest) GetValidity() string {
	return request.validity
}

// SetCallbackURL set the callback url this url will get called when there is a status update
func (request *SmsRequest) SetCallbackURL(URL *url.URL) {
	request.callbackURL = URL
}

// GetCallbackURL get the callback url
func (request SmsRequest) GetCallbackURL() *url.URL {
	return request.callbackURL
}

// SetDCS sets the DCS of an SMS Request [0 = GSM-7, 8 = unicode] will automatically be set if not explicitly set
// allowed values 0-255, see https://en.wikipedia.org/wiki/GSM_03.40#Data_Coding_Scheme
func (request *SmsRequest) SetDCS(dcs int) {
	request.dcs = dcs
}

// GetDcs get the DCS set for the SMS Request
func (request SmsRequest) GetDcs() int {
	return request.dcs
}

// SetUdh sets the Usage Data Header (used for concat messages), can only consist of hexadecimal characters,
// length of this parameter divided by 2 is subtracted form the max message size,
// see https://en.wikipedia.org/wiki/User_Data_Header
func (request *SmsRequest) SetUdh(udh *string) {
	request.udh = udh
}

// GetUdh get the Usage Data Header
func (request SmsRequest) GetUdh() *string {
	return request.udh
}

// Submit the message
func (request *SmsRequest) Submit() (*SmsResponses, error) {
	response := &SmsResponses{}

	var apiURL *url.URL
	var err error

	if request.submitType == SmsSubmitTypeSimple {
		apiURL, err = GetURLFor("sms/submitsimple")
	} else {
		apiURL, err = GetURLFor("sms/submit")
	}

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

// SmsResponse structure response from status
type SmsResponse struct {
	applicationTag         string
	body                   []byte
	callbackURL            *url.URL
	createDateTime         time.Time
	dcs                    int
	messageID              string
	networkCode            *int
	pid                    *string
	reasonCode             *int
	recipient              Recipient
	resultTimestamp        *string
	resultType             int
	salesPrice             *float32
	salesPriceCurrencyCode *string
	scheduledDelivery      *string
	sender                 string
	senderNpi              int
	senderTon              int
	statusMsg              string
	statusCode             SmsStatusCode
	tag                    *string
	udh                    *string
	validity               int
	validUntilDateTime     time.Time
	links                  HATEOASLinks
}

type jsonSmsResponse struct {
	ApplicationTag         string        `json:"applicationTag,omitempty"`
	Body                   string        `json:"body"`
	CallbackURL            *url.URL      `json:"callbackUrl,omitempty"`
	CreateDateTime         time.Time     `json:"createdDateTime,omitempty"`
	Dcs                    int           `json:"dcs,omitempty"`
	MessageID              string        `json:"messageId"`
	NetworkCode            *int          `json:"networkCode,omitempty"`
	Pid                    *string       `json:"pid,omitempty"`
	ReasonCode             *int          `json:"reasonCode,omitempty"`
	Recipient              Recipient     `json:"recipient"`
	ResultTimestamp        *string       `json:"resultTimestamp,omitempty"`
	ResultType             int           `json:"resultType,omitempty"`
	SalesPrice             *float32      `json:"salesPrice,omitempty"`
	SalesPriceCurrencyCode *string       `json:"salesPriceCurrencyCode,omitempty"`
	ScheduledDelivery      *string       `json:"scheduledDelivery,omitempty"`
	Sender                 string        `json:"sender"`
	SenderNpi              int           `json:"senderNpi,omitempty"`
	SenderTon              int           `json:"senderTon,omitempty"`
	StatusMsg              string        `json:"status,omitempty"`
	StatusCode             SmsStatusCode `json:"statusCode,omitempty"`
	Tag                    *string       `json:"tag,omitempty"`
	Udh                    *string       `json:"udh,omitempty"`
	Validity               int           `json:"validity,omitempty"`
	ValidUntilDateTime     time.Time     `json:"validUntilDateTime,omitempty"`
	Links                  HATEOASLinks  `json:"_links"`
}

// UnmarshalJSON convert the json sms response to a jsonResponse struct
func (response *SmsResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonSmsResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *SmsResponse) copyFrom(j *jsonSmsResponse) error {
	var err error // default err is nil

	response.applicationTag = j.ApplicationTag
	response.callbackURL = j.CallbackURL
	response.createDateTime = j.CreateDateTime
	response.dcs = j.Dcs
	response.messageID = j.MessageID
	response.networkCode = j.NetworkCode
	response.pid = j.Pid
	response.reasonCode = j.ReasonCode
	response.recipient = j.Recipient
	response.resultTimestamp = j.ResultTimestamp
	response.tag = j.Tag
	response.resultType = j.ResultType
	response.salesPrice = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.scheduledDelivery = j.ScheduledDelivery
	response.sender = j.Sender
	response.senderNpi = j.SenderNpi
	response.senderTon = j.SenderTon
	response.statusMsg = j.StatusMsg
	response.statusCode = j.StatusCode
	response.udh = j.Udh
	response.validity = j.Validity
	response.validUntilDateTime = j.ValidUntilDateTime
	response.links = j.Links

	if response.IsBinary() {
		response.body, err = hex.DecodeString(j.Body)
	} else {
		response.body = []byte(j.Body)
	}

	return err
}

// IsBinary returns true if the send sms was binary
func (response SmsResponse) IsBinary() bool {
	return isDcsBinary(response.dcs)
}

// GetApplicationTag gets the application tag for the sms response
func (response SmsResponse) GetApplicationTag() string {
	return response.applicationTag
}

// GetBodyAsString get the body as string
func (response SmsResponse) GetBodyAsString() string {
	return string(response.body)
}

// GetBodyAsByteArr get the body as []byte
func (response SmsResponse) GetBodyAsByteArr() []byte {
	return response.body
}

// GetCallbackURL returns the callback set on the response
func (response SmsResponse) GetCallbackURL() *url.URL {
	return response.callbackURL
}

// GetCreateDateTime returns the time that the sms was sent
func (response SmsResponse) GetCreateDateTime() time.Time {
	return response.createDateTime
}

// GetDcs gets the dcs set by the request stuct or the automatic dcs set by the server
func (response SmsResponse) GetDcs() int {
	return response.dcs
}

// GetMessageID gets the message id set by the server
func (response SmsResponse) GetMessageID() string {
	return response.messageID
}

// GetNetworkCode gets the network code if set by the server
func (response SmsResponse) GetNetworkCode() *int {
	return response.networkCode
}

// GetPid returns the Protocol Identifier set by the server
func (response SmsResponse) GetPid() *string {
	return response.pid
}

// GetReasonCode set by the server
func (response SmsResponse) GetReasonCode() *int {
	return response.reasonCode
}

// GetRecipient returns the recipient that this sms was sent to
func (response SmsResponse) GetRecipient() Recipient {
	return response.recipient
}

// GetResultTimestamp returns the timestamp of the result
func (response SmsResponse) GetResultTimestamp() *string {
	return response.resultTimestamp
}

// GetResultType returns the type of the result
func (response SmsResponse) GetResultType() int {
	return response.resultType
}

// GetSalesPrice returns the sales price of the sent sms
func (response SmsResponse) GetSalesPrice() *float32 {
	return response.salesPrice
}

// GetSalesPriceCurrencyCode returns the currency of the sales price
func (response SmsResponse) GetSalesPriceCurrencyCode() *string {
	return response.salesPriceCurrencyCode
}

// GetScheduledDelivery returns the scheduled delivery string if set
func (response SmsResponse) GetScheduledDelivery() *string {
	return response.scheduledDelivery
}

// GetSender returns the sender of the sms
func (response SmsResponse) GetSender() string {
	return response.sender
}

// GetSenderNpi returns the senderNpi see: https://www.twizo.com/tutorials/#sender
func (response SmsResponse) GetSenderNpi() int {
	return response.senderNpi
}

// GetSenderTon returns the senderTon see: https://www.twizo.com/tutorials/#sender
func (response SmsResponse) GetSenderTon() int {
	return response.senderTon
}

// GetStatusMsg returns the status message set by the server
func (response SmsResponse) GetStatusMsg() string {
	return response.statusMsg
}

// GetStatusCode returns the status code set by the server
func (response SmsResponse) GetStatusCode() SmsStatusCode {
	return response.statusCode
}

// GetTag gets the tag of the sms response if set
func (response SmsResponse) GetTag() *string {
	return response.tag
}

// GetUdh returns the udh of the sent sms set: https://en.wikipedia.org/wiki/User_Data_Header
func (response SmsResponse) GetUdh() *string {
	return response.udh
}

// GetValidity returns the validity set on the sms sent
func (response SmsResponse) GetValidity() int {
	return response.validity
}

// GetValidUntilDateTime returns the time.Time that the sms was valid for
func (response SmsResponse) GetValidUntilDateTime() time.Time {
	return response.validUntilDateTime
}

// Status gets the status of one message
func (response *SmsResponse) Status() error {
	newResponse := &SmsResponse{}

	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		&response.links.Self.Href,
		nil,
		http.StatusOK,
		newResponse,
	)

	if err == nil {
		*response = *newResponse
	}

	return err
}

type smsResponseEmbedded struct {
	Items *[]SmsResponse `json:"items"`
}

type jsonSmsResponses struct {
	Links    HATEOASLinks        `json:"_links"`
	Embedded smsResponseEmbedded `json:"_embedded"`
	Total    int                 `json:"total_items"`
}

// SmsResponses response contains multiple SmsMessages
type SmsResponses struct {
	Responses *[]SmsResponse
}

// UnmarshalJSON unmarshals the responses to JSON
func (r *SmsResponses) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonSmsResponses{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	r.Responses = jsonResponse.Embedded.Items

	return nil
}

// Status requests the status of a response
func (r *SmsResponses) Status() error {
	newItems := []SmsResponse{}
	for _, item := range *r.Responses {
		err := item.Status()
		if err != nil {
			return err
		}
		newItems = append(newItems, item)
	}
	// all went well update Items
	r.Responses = &newItems
	return nil
}

// GetItems retrieve all messages for the response
func (r *SmsResponses) GetItems() []SmsResponse {
	return *r.Responses
}

type smsPollEmbedded struct {
	Messages *[]SmsResponse `json:"messages"` // strange this was called items in smsResponseEmbedded
}

// SmsPollResults struct
type SmsPollResults struct {
	pollResults                  // an anonymous field of type pollResult
	Embedded    *smsPollEmbedded `json:"_embedded"`
}

func (p SmsPollResults) getURL() (*url.URL, error) {
	return GetURLFor("sms/poll")
}

// GetItems get the embedded items of a poll result
func (p SmsPollResults) GetItems() []SmsResponse {
	return *p.Embedded.Messages
}

// NewSmsPoll creates a new smspollresult
func NewSmsPoll() *SmsPollResults {
	params := &SmsPollResults{}
	return params
}

// SmsPollStatus [todo: refactor]
func SmsPollStatus() (*SmsPollResults, error) {
	request := NewSmsPoll()
	err := request.Status()

	return request, err
}

// Status get the status of a sms poll result
func (p *SmsPollResults) Status() error {
	apiURL, err := p.getURL()
	if err != nil {
		return err
	}
	resp := NewSmsPoll()
	err = p.status(apiURL, resp)
	if err == nil {
		*p = *resp
	}
	return err
}

// NewSmsRequest creates a new smsrequest struct
func NewSmsRequest(recipients []Recipient, body interface{}, sender string) (*SmsRequest, error) {
	params := &SmsRequest{
		recipients: recipients,
		submitType: SmsSubmitTypeSimple,
	}
	if err := params.SetSender(sender); err != nil {
		return nil, err
	}
	if err := params.SetBody(body); err != nil {
		return nil, err
	}

	return params, nil
}

// SmsStatus retrieves the status of a message by ID
func SmsStatus(messageID string) (*SmsResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("sms/submit/%s", url.PathEscape(messageID)))
	if err != nil {
		return nil, err
	}

	smsResponse := &SmsResponse{
		messageID: messageID,
		links:     createSelfLinks(apiURL),
	}

	err = smsResponse.Status()
	// were we able to find it ?
	if apiError, ok := err.(APIError); ok {
		if apiError.Status() == http.StatusNotFound {
			return nil, nil
		}
	}

	return smsResponse, err
}

// SmsSubmit submits a message to recipients
func SmsSubmit(recipients interface{}, body interface{}, sender string) (*SmsResponses, error) {
	r, err := convertRecipients(recipients)
	if err != nil {
		return nil, err
	}
	sms, err := NewSmsRequest(r, body, sender)
	if err != nil {
		return nil, err
	}

	return sms.Submit()
}
