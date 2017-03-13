package twizo

import (
	"fmt"
	"net/http"
)

type NumberLookupStatusCode int

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


// Create a new verificationParam using a recipient (the only required var)
func NewNumberLookupRequest(numbers []Recipient) *NumberLookupRequest {
	params := &NumberLookupRequest{
		numbers: numbers,
	}
	return params
}

func NumberLookupSubmit(numbers interface{}) (*NumberLookupResponses, error) {
	r, err := convertRecipients(numbers)
	if err != nil {
		return nil, err
	}
	return NewNumberLookupRequest(r).Submit()
}

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
