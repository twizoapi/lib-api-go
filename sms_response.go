package twizo

import (
	"net/url"
	"time"
	"net/http"
	"encoding/json"
	"encoding/hex"
)

// SmsMessage stucture response from status
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
	statusCode             smsStatusCode
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
	Pid                    *string        `json:"pid,omitempty"`
	ReasonCode             *int           `json:"reasonCode,omitempty"`
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
	StatusCode             smsStatusCode `json:"statusCode,omitempty"`
	Tag                    *string       `json:"tag,omitempty"`
	Udh                    *string       `json:"udh,omitempty"`
	Validity               int           `json:"validity,omitempty"`
	ValidUntilDateTime     time.Time     `json:"validUntilDateTime,omitempty"`
	Links                  HATEOASLinks  `json:"_links"`
}

func (r *SmsResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonSmsResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	r.copyFrom(jsonResponse)

	return nil
}

func (response *SmsResponse) copyFrom (j *jsonSmsResponse) (error) {
	var err error  // default err is nil

	response.applicationTag 	= j.ApplicationTag
	response.callbackURL 		= j.CallbackURL
	response.createDateTime 	= j.CreateDateTime
	response.dcs			= j.Dcs
	response.messageID		= j.MessageID
	response.networkCode            = j.NetworkCode
	response.pid			= j.Pid
	response.reasonCode		= j.ReasonCode
	response.recipient              = j.Recipient
	response.resultTimestamp	= j.ResultTimestamp
	response.tag                    = j.Tag
	response.resultType             = j.ResultType
	response.salesPrice             = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.scheduledDelivery      = j.ScheduledDelivery
	response.sender                 = j.Sender
	response.senderNpi              = j.SenderNpi
	response.senderTon              = j.SenderTon
	response.statusMsg              = j.StatusMsg
	response.statusCode             = j.StatusCode
	response.udh                    = j.Udh
	response.validity               = j.Validity
	response.validUntilDateTime     = j.ValidUntilDateTime
	response.links                  = j.Links

	if response.IsBinary() {
		response.body, err = hex.DecodeString(j.Body)
	} else {
		response.body = []byte(j.Body)
	}

	return err
}

func (response SmsResponse) IsBinary() (bool) {
	return isDcsBinary(response.dcs)
}

// GetApplicationTag
func (response SmsResponse) GetApplicationTag() (string) {
	return response.applicationTag
}

// GetBodyAsString get the body as string
func (response SmsResponse) GetBodyAsString() (string) {
	return string(response.body)
}

// GetBodyAsByteArr get the body as []byte
func (response SmsResponse) GetBodyAsByteArr() ([]byte) {
	return response.body
}

func (response SmsResponse) GetCallbackURL() (*url.URL) {
	return response.callbackURL
}

func (response SmsResponse) GetCreateDateTime() (time.Time) {
	return response.createDateTime
}

func (response SmsResponse) GetDcs() (int) {
	return response.dcs
}

func (response SmsResponse) GetMessageID() (string) {
	return response.messageID
}

func (response SmsResponse) GetNetworkCode() (*int) {
	return response.networkCode
}

func (response SmsResponse) GetPid() (*string) {
	return response.pid
}

func (response SmsResponse) GetReasonCode() (*int) {
	return response.reasonCode
}

func (response SmsResponse) GetRecipient() (Recipient) {
	return response.recipient
}

func (response SmsResponse) GetResultTimestamp() (*string) {
	return response.resultTimestamp
}

func (response SmsResponse) GetResultType() (int) {
	return response.resultType
}

func (response SmsResponse) GetSalesPrice() (*float32) {
	return response.salesPrice
}

func (response SmsResponse) GetSalesPriceCurrencyCode() (*string) {
	return response.salesPriceCurrencyCode
}

func (response SmsResponse) GetScheduledDelivery() (*string) {
	return response.scheduledDelivery
}

func (response SmsResponse) GetSender() (string) {
	return response.sender
}

func (response SmsResponse) GetSenderNpi() (int) {
	return response.senderNpi
}

func (response SmsResponse) GetSenderTon() (int) {
	return response.senderTon
}

func (response SmsResponse) GetStatusMsg() (string) {
	return response.statusMsg
}

func (response SmsResponse) GetStatusCode() (smsStatusCode) {
	return response.statusCode
}

// GetTag
func (response SmsResponse) GetTag() (*string) {
	return response.tag
}

func (response SmsResponse) GetUdh() (*string) {
	return response.udh
}

func (response SmsResponse) GetValidity() (int) {
	return response.validity
}

func (response SmsResponse) GetValidUntilDateTime() (time.Time) {
	return response.validUntilDateTime
}

// Status gets the status of one message
func (response *SmsResponse) Status() (error) {
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

// / SmsResponse response containting multiple SmsMessages
type SmsResponses struct {
	Responses *[]SmsResponse
}

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
