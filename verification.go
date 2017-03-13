package twizo

import (
	"fmt"
)

type VerificationStatusCode int
type VerificationType string
type VerificationTokenType string

type VerificationTokenErrorCode int

const (
	// Verification Type
	VerificationTypeSms  VerificationType = "sms"
	VerificationTypeCall VerificationType = "call"

	// Verification Token Type
	VerificationTokenTypeDefault VerificationTokenType = ""
	VerificationTokenTypeNumeric VerificationTokenType = "numeric"
	VerificationTokenTypeAlpha   VerificationTokenType = "alphanumeric"

	// Verification Status
	VerificationTokenUnknown         VerificationStatusCode = 0
	VerificationTokenSuccess         VerificationStatusCode = 1
	VerificationTokenAlreadyVerified VerificationStatusCode = 101
	VerificationTokenExpired         VerificationStatusCode = 102
	VerificationTokenInvalid         VerificationStatusCode = 103
	VerificationTokenFailed          VerificationStatusCode = 104
)

// Create a new verificationParam using a recipient (the only required var)
func NewVerificationRequest(recipient Recipient) *VerificationRequest {
	params := &VerificationRequest{recipient: recipient}
	return params
}


// Create a verificationRequest from verificationParams and submit it
func VerificationSubmit(recipient interface{}) (*VerificationResponse, error) {
	r, err := convertRecipients(recipient)
	if err != nil {
		return nil, err
	}
	if len(r) != 1 {
		return nil, fmt.Errorf("Need exactly one [recipient] for VerificationSubmit got [%d]", len(r))
	}

	return NewVerificationRequest(r[0]).Submit()
}

// Retrieve the status of the validation using the messageId
func VerificationStatus(messageID string) (*VerificationResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("verification/submit/%s", messageID))
	if err != nil {
		return nil, err
	}

	request := &VerificationResponse{messageID: messageID, links: createSelfLinks(apiURL)}
	err = request.Status()
	return request, err
}

// Validate the result of the validation request using the messageId and the token
func VerificationVerify(messageID string, token string) (*VerificationResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("verification/submit/%s", messageID))
	if err != nil {
		return nil, err
	}

	request := &VerificationResponse{messageID: messageID, links: createSelfLinks(apiURL)}
	err = request.Verify(token)
	return request, err
}
