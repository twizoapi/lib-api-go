package twizo

import (
	"fmt"
	"net/http"
)

type smsSubmitType int
type smsStatusCode int

const (
	// SmsSubmitTypeSimple (default) use simple submit
	// (auto body splitting auto udh / dcs setting)
	SmsSubmitTypeSimple             smsSubmitType = 0

	// SmsSubmitTypeAdvanced will not autosplit and can
	// also send binary messages
	SmsSubmitTypeAdvanced        	smsSubmitType = 1

	// SmsStatusCodeNoStatus no status for message
	SmsStatusCodeNoStatus    	smsStatusCode = 0

	// SmsStatusCodeDelivered message delivered
	SmsStatusCodeDelivered   	smsStatusCode = 1

	// SmsStatusCodeRejected message rejected
	SmsStatusCodeRejected    	smsStatusCode = 2

	// SmsStatusCodeExpired message expired
	SmsStatusCodeExpired     	smsStatusCode = 3

	// SmsStatusCodeEnroute message enroute
	SmsStatusCodeEnroute     	smsStatusCode = 4

	// SmsStatusCodeBuffered message buffered
	SmsStatusCodeBuffered    	smsStatusCode = 5

	// SmsStatusCodeAccepted message accepted
	SmsStatusCodeAccepted    	smsStatusCode = 6

	// SmsStatusCodeUndelivered message undelivered
	SmsStatusCodeUndelivered 	smsStatusCode = 7

	// SmsStatusCodeDeleted message deleted
	SmsStatusCodeDeleted     	smsStatusCode = 8

	// SmsStatusCodeUnknown message status unkown
	SmsStatusCodeUnknown     	smsStatusCode = 9
)

// NewSmsRequest creates a new smsrequest struct
func NewSmsRequest(recipients []Recipient, body interface{}, sender string) *SmsRequest {
	params := &SmsRequest{
		recipients: recipients,
		submitType: SmsSubmitTypeSimple,
	}
	params.SetSender(sender)
	params.SetBody(body)

	return params
}

// SmsStatus retrieves the status of a message by ID
func SmsStatus(messageID string) (*SmsResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("sms/submit/%s", messageID))
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
		if apiError.Status == http.StatusNotFound {
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
	return NewSmsRequest(r, body, sender).Submit()
}
