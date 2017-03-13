package twizo

import (
	"net/url"
	"time"
	"net/http"
	"encoding/json"
)

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
	CallbackURL            *url.URL               `json:"callbackUrl,omitempty"`
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

func (r *NumberLookupResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonNumberLookupResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	r.copyFrom(jsonResponse)

	return nil
}

func (response *NumberLookupResponse) copyFrom (j *jsonNumberLookupResponse) (error) {
	var err error  // default err is nil

	response.applicationTag 	= j.ApplicationTag
	response.callbackURL 		= j.CallbackURL
	response.createDateTime 	= j.CreateDateTime
	response.imsi			= j.Imsi
	response.ported                 = j.Ported
	response.roaming                = j.Roaming
	response.messageID		= j.MessageID
	response.msc			= j.Msc
	response.networkCode            = j.NetworkCode
	response.number			= j.Number
	response.operator               = j.Operator
	response.reasonCode             = j.ReasonCode
	response.resultTimestamp	= j.ResultTimestamp
	response.resultType             = j.ResultType
	response.salesPrice             = j.SalesPrice
	response.salesPriceCurrencyCode = j.SalesPriceCurrencyCode
	response.statusMsg              = j.StatusMsg
	response.statusCode             = j.StatusCode
	response.tag                    = j.Tag
	response.validity               = j.Validity
	response.validUntilDateTime     = j.ValidUntilDateTime
	response.links                  = j.Links

	return err
}

// Status
func (n *NumberLookupResponse) Status() (error) {
	response := &NumberLookupResponse{}

	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		&n.links.Self.Href,
		nil,
		http.StatusOK,
		response,
	)

	if err == nil {
		// no error use response to override ourselves
		*n = *response
	}

	return err
}

func (n NumberLookupResponse) GetApplicationTag() (string) {
	return n.applicationTag
}

func (n NumberLookupResponse) GetCallbackURL() (*url.URL) {
	return n.callbackURL
}

func (n NumberLookupResponse) GetCreateDateTime() (time.Time) {
	return n.createDateTime
}

func (n NumberLookupResponse) GetImsi() (*string) {
	return n.imsi
}

func (n NumberLookupResponse) GetPorted() (string) {
	return n.ported
}

func (n NumberLookupResponse) GetRoaming() (string) {
	return n.roaming
}

func (n NumberLookupResponse) GetMessageID() (string) {
	return n.messageID
}

func (n NumberLookupResponse) GetMsc() (*string) {
	return n.msc
}

func (n NumberLookupResponse) GetNetworkCode() (*int) {
	return n.networkCode
}

func (n NumberLookupResponse) GetNumber() (string) {
	return n.number
}

func (n NumberLookupResponse) GetOperator() (*string) {
	return n.operator
}

func (n NumberLookupResponse) GetReasonCode() (*int) {
	return n.reasonCode
}

func (n NumberLookupResponse) GetResultType() (int) {
	return n.resultType
}

func (n NumberLookupResponse) GetSalesPrice() (*float32) {
	return n.salesPrice
}

func (n NumberLookupResponse) GetSalesPriceCurrencyCode() (*string) {
	return n.salesPriceCurrencyCode
}

func (n NumberLookupResponse) GetStatusMsg() (string) {
	return n.statusMsg
}

func (n NumberLookupResponse) GetStatusCode() (NumberLookupStatusCode) {
	return n.statusCode
}

func (n NumberLookupResponse) GetTag() (*string) {
	return n.tag
}

func (n NumberLookupResponse) GetResultTimeStamp() (string) {
	return n.resultTimestamp
}

func (n NumberLookupResponse) GetValidity() (int) {
	return n.validity
}

func (n NumberLookupResponse) GetValidUntillDateTime() (time.Time) {
	return n.validUntilDateTime
}

type numberLookupEmbedded struct {
	Items *[]NumberLookupResponse `json:"items"`
}

//
// Request status of all embedded objects
//
func (v *numberLookupEmbedded) Status() error {
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

type NumberLookupResponses struct {
	Responses *[]NumberLookupResponse
}

type jsonNumberLookupResponses struct {
	Links    HATEOASLinks         `json:"_links"`
	Embedded numberLookupEmbedded `json:"_embedded"`
	Total    int                  `json:"total_items"`
}

func (r *NumberLookupResponses) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonNumberLookupResponses{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	r.Responses = jsonResponse.Embedded.Items

	return nil
}

//
// Get all messages
//
func (r NumberLookupResponses) GetItems() []NumberLookupResponse {
	return *r.Responses
}

//
// Request status of all embedded objects
//
func (r *NumberLookupResponses) Status() error {
	newItems := []NumberLookupResponse{}
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
