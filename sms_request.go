package twizo

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
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
