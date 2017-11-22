package twizo

import (
	"fmt"
	"net/http"
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
	apiURL, err := GetURLFor(fmt.Sprintf("numberlookup/submit/%s", messageID))
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
		if apiError.Status == http.StatusNotFound {
			return nil, nil
		}
	}

	return numberLookupResponse, err
}
