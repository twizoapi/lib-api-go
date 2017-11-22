package twizo

import (
	"fmt"
)
// VerificationStatusCode contains the verification status code
type VerificationStatusCode int

// VerificationType describes the verification type
type VerificationType string

// VerificationTokenType describes the verification token type
type VerificationTokenType string

// VerificationTokenErrorCode describes the token error code
type VerificationTokenErrorCode int

// All verification types supported
const (
	VerificationTypeSms  VerificationType = "sms"
	VerificationTypeCall VerificationType = "call"
)

// All verification Token types
const (
	VerificationTokenTypeDefault VerificationTokenType = ""
	VerificationTokenTypeNumeric VerificationTokenType = "numeric"
	VerificationTokenTypeAlpha   VerificationTokenType = "alphanumeric"
)

// Result of verification tokens
const (
	VerificationTokenUnknown         VerificationStatusCode = 0
	VerificationTokenSuccess         VerificationStatusCode = 1
	VerificationTokenAlreadyVerified VerificationStatusCode = 101
	VerificationTokenExpired         VerificationStatusCode = 102
	VerificationTokenInvalid         VerificationStatusCode = 103
	VerificationTokenFailed          VerificationStatusCode = 104
)

// NewVerificationRequest creates a new verificationParam using a recipient (the only required var)
func NewVerificationRequest(recipient Recipient) *VerificationRequest {
	params := &VerificationRequest{recipient: recipient}
	return params
}

// VerificationSubmit creates a verificationRequest from verificationParams and submits it
func VerificationSubmit(recipient interface{}) (*VerificationResponse, error) {
	r, err := convertRecipients(recipient)
	if err != nil {
		return nil, err
	}
	if len(r) != 1 {
		return nil, fmt.Errorf("need exactly one [recipient] for VerificationSubmit got [%d]", len(r))
	}

	return NewVerificationRequest(r[0]).Submit()
}

// VerificationStatus retrieves the status of the validation using the messageId
func VerificationStatus(messageID string) (*VerificationResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("verification/submit/%s", messageID))
	if err != nil {
		return nil, err
	}

	request := &VerificationResponse{messageID: messageID, links: createSelfLinks(apiURL)}
	err = request.Status()
	return request, err
}

// VerificationVerify validates the result of the validation request using the messageId and the token
func VerificationVerify(messageID string, token string) (*VerificationResponse, error) {
	apiURL, err := GetURLFor(fmt.Sprintf("verification/submit/%s", messageID))
	if err != nil {
		return nil, err
	}

	request := &VerificationResponse{messageID: messageID, links: createSelfLinks(apiURL)}
	err = request.Verify(token)
	return request, err
}
