package twizo

import (
	"fmt"
	"net/http"
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
	sms, err := NewSmsRequest(r, body, sender)
	if err != nil {
		return nil, err
	}

	return sms.Submit()
}
