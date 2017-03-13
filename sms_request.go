package twizo

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"encoding/hex"
	"reflect"
)

// SmsRequest is used to send an sms request
type SmsRequest struct {
	recipients        []Recipient
	body              []byte         // 10 x 160
	submitType        smsSubmitType  // type simple (default) or advanced (not sent to server)
	sender            string
	senderTon         int		 // posivive small int
	senderNpi         int        	 // positive small int
	pid               int
	scheduledDelivery string
	tag               string
	validity          string         // duration.Seconds
	resultType        ResultType
	callbackURL       *url.URL       // only relevant for SmsResultTypeCallback | SmsResultTYpeCallbackPollling
	dcs               int            // relevant for advanced, 0-255
	udh               *string        // relevant for advanced, hex
}

type jsonSmsRequest struct {
	Recipients        []Recipient    `json:"recipients"`
	Body              string         `json:"body"`
	Sender            string         `json:"sender"`
	SenderTon         int            `json:"senderTon,omitempty"`
	SenderNpi         int            `json:"senderNpi,omitempty"`
	Pid               int            `json:"pid,omitempty"`
	ScheduledDelivery string         `json:"scheduledDelivery,omitempty"`
	Tag               string         `json:"tag,omitempty"`
	Validity          string         `json:"validity,omitempty"`
	ResultType        ResultType     `json:"resultType"`
	CallbackURL       *url.URL       `json:"callbackUrl,omitempty"`
	Dcs               int            `json:"dcs,omitempty"`
	Udh               *string        `json:"udh,omitempty"`
}

// MarshalJSON is used to convert SmsRequest to json
func (request *SmsRequest) MarshalJSON() ([]byte, error) {
	jsonSmsRequest := jsonSmsRequest{
		Recipients 	  : request.recipients,
		Body 		  : string(request.body),
		Sender 		  : request.sender,
		SenderTon         : request.senderTon,
		SenderNpi         : request.senderNpi,
		Pid               : request.pid,
		ScheduledDelivery : request.scheduledDelivery,
		Validity          : request.validity,
		ResultType        : request.resultType,
		Tag               : request.tag,
	}

	// set the callback url if we need one, it still might be empty
	if request.resultType == ResultTypeCallback || request.resultType == ResultTypeCallbackPolling {
		jsonSmsRequest.CallbackURL = request.callbackURL
	}

	if request.submitType == SmsSubmitTypeAdvanced {
		// if request is advanced, send dcs and udh too
		jsonSmsRequest.Dcs = request.dcs
		jsonSmsRequest.Udh = request.udh

		// if the request is binary we hex encode the body
		if request.IsBinary() {
			jsonSmsRequest.Body = hex.EncodeToString(request.body)
		}
	}

	return json.Marshal(jsonSmsRequest)
}

// IsBinary returns if message is binary or not
func (request SmsRequest) IsBinary() (bool) {
	return isDcsBinary(request.dcs)
}

// SetSender
func (r *SmsRequest) SetSender(sender string) (error) {
	r.sender = sender

	return nil
}

// GetSender
func (r SmsRequest) GetSender() (string, error) {
	return r.sender, nil
}

// SetTag
func (r *SmsRequest) SetTag(tag string) (error) {
	r.tag = tag

	return nil
}

// GetTag
func (r SmsRequest) GetTag() (string, error) {
	return r.tag, nil
}

// SetBody set the body
func (request *SmsRequest) SetBody(body interface{}) (error) {
	switch t := body.(type) {
	case string:
		return request.SetBodyFromString(t)
	case []byte:
		return request.SetBodyFromByteArr(t)
	default:
		return fmt.Errorf("SetBody expects string or []byte, got [%s]", reflect.TypeOf(body))
	}

	return nil
}

// SetBodyFromString set the body as string
func (request *SmsRequest) SetBodyFromString(body string) (error) {
	request.body = []byte(body)
	return nil
}

// SetBodyFromByteArr set the body as string
func (request *SmsRequest) SetBodyFromByteArr(body []byte) (error) {
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
func (request *SmsRequest) SetRecipients(recipients []Recipient) (error) {
	request.recipients = recipients

	return nil
}

// SetResultType set the resultType
func (request *SmsRequest) SetResultType(resultType ResultType) (error) {
	request.resultType = resultType

	return nil
}
// GetResultType get the resultType
func (request SmsRequest) GetResultType() (ResultType, error) {
	return request.resultType, nil
}

func (request *SmsRequest) SetSenderTon(ton int) {
	request.senderTon = ton
}

func (request SmsRequest) GetSenderTon() (int) {
	return request.senderTon
}

func (request *SmsRequest) SetSenderNpi(npi int) {
	request.senderNpi = npi
}

func (request SmsRequest) GetSenderNpi() (int) {
	return request.senderNpi
}

func (request *SmsRequest) SetPid(pid int) {
	request.pid = pid
}

func (request SmsRequest) GetPid() (int) {
	return request.pid
}

func (request *SmsRequest) SetScheduledDelivery(sd string) {
	request.scheduledDelivery = sd
}

func (request SmsRequest) GetScheduledDelivery() (string) {
	return request.scheduledDelivery
}

func (request *SmsRequest) SetValidity(validity string) {
	request.validity = validity
}

func (request SmsRequest) GetValidity() (string) {
	return request.validity
}

func (request *SmsRequest) SetCallbackURL(URL *url.URL) {
	request.callbackURL = URL
}

func (request SmsRequest) GetCallbackURL() (*url.URL) {
	return request.callbackURL
}

func (request *SmsRequest) SetDCS(dcs int) {
	request.dcs = dcs
}

func (request SmsRequest) GetDcs() (int) {
	return request.dcs
}

func (request *SmsRequest) SetUdh(udh *string) {
	request.udh = udh
}

func (request SmsRequest) GetUdh() (*string) {
	return request.udh
}

// Submit the message
func (request *SmsRequest) Submit() (*SmsResponses, error) {
	response := &SmsResponses{}

	var apiURL *url.URL
	var err    error

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
