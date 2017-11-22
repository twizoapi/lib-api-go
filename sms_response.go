package twizo

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

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
