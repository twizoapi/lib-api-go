package twizo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

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
	response.callbackURL = j.CallbackURL
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

type numberLookupEmbedded struct {
	Items *[]NumberLookupResponse `json:"items"`
}

// Status requests the status of all embedded numberlookups
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

// NumberLookupResponses struct
type NumberLookupResponses struct {
	Responses *[]NumberLookupResponse
}

type jsonNumberLookupResponses struct {
	Links    HATEOASLinks         `json:"_links"`
	Embedded numberLookupEmbedded `json:"_embedded"`
	Total    int                  `json:"total_items"`
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
