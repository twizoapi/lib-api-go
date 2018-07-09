package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	"encoding/json"
	"testing"

	"net/http"

	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"
)

func init() {
	twizo.APIKey = TestAPIKey
	twizo.RegionCurrent = TestRegion
}

func TestAPIErrorInvalidJsonResponse(t *testing.T) {
	apiError := &twizo.APIError{}
	err := apiError.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestAPIValidationErrorInvalidJsonResponse(t *testing.T) {
	apiError := &twizo.APIValidationError{}
	err := apiError.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	} else {
	}
}

func TestAPIErrorInvalidResponse(t *testing.T) {
	cannedResponse := `{
		  "type":"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
		  "title":"Unprocessable Entity",
		  "status":422,
		  "detail":"Failed Validation"
         }`

	apiError := &twizo.APIError{}
	err := json.Unmarshal([]byte(cannedResponse), apiError)
	if err != nil {
		t.Fatal(err)
	}
	if apiError.Detail() != "Failed Validation" {
		t.Fatalf(
			"Invalid detail expecting [Failed Validation] got [%v]",
			apiError.Detail(),
		)
	}
}

func TestAPIErrorVerificationOneErrorResponse(t *testing.T) {
	cannedResponse := `{
		"validation_messages": {
			"allowedTypes": {
				"1": {
					"value":"",
					"validation_errors": {
						"notInArray":"Invalid type '' specified; only 'sms', 'call', ` +
		`        'backupcode', 'biovoice', 'totpCreate', 'push', 'line' or ` +
		`        'telegram' allowed"
					}
				}
			}
		},
		"type":"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
		"title":"Unprocessable Entity",
		"status":422,
		"detail":"Failed Validation"
	 }`

	apiError := &twizo.APIValidationError{}
	err := json.Unmarshal([]byte(cannedResponse), apiError)
	if err != nil {
		t.Fatal(err)
	}
	if apiError.Detail() != "Failed Validation" {
		t.Fatalf(
			"Invalid detail expecting [Failed Validation] got [%v]",
			apiError.Detail(),
		)
	}
	if apiError.Title() != "Unprocessable Entity" {
		t.Fatalf(
			"Invalid detail expecting [Failed Validation] got [%v]",
			apiError.Title(),
		)
	}
	if apiError.Status() != http.StatusUnprocessableEntity {
		t.Fatalf(
			"Invalid detail expecting [%d] got [%d]",
			http.StatusUnprocessableEntity,
			apiError.Status(),
		)
	}
}

func TestAPIErrorVerificationTwoErrorsResponse(t *testing.T) {
	cannedResponse := `{
		"validation_messages": {
			"recipients": {
				"1": {
					"value":"g",
					"validation_errors": {
						"stringLengthTooShort":"The input is less than 8 characters long",
						"notDigits":"The input must contain only digits"
					}
				}
			}
		},
		"type":"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
		"title":"Unprocessable Entity",
		"status":422,
		"detail":"Failed Validation"
	}`
	apiError := &twizo.APIValidationError{}
	err := json.Unmarshal([]byte(cannedResponse), apiError)
	if err != nil {
		t.Fatal(err)
	}
	if apiError.Detail() != "Failed Validation" {
		t.Fatalf(
			"Invalid detail expecting [Failed Validation] got [%v]",
			apiError.Detail(),
		)
	}
}

func TestAPIErrorVerificationFieldNotAllowedResponse(t *testing.T) {
	cannedResponse := `{
		"validation_messages": {
			"-": {
				"invalidFields":"The following field(s) are not allowed: 'test'"
			},
			"allowedTypes": {
				"1": {
					"value":"",
					"validation_errors": {
						"notInArray":"Invalid type '' specified; only 'sms', 'call', ` +
		`        'backupcode', 'biovoice', 'totpCreate', 'push', 'line' or ` +
		`        'telegram' allowed"
					}
				}
			}
		},
		"type":"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
		"title":"Unprocessable Entity",
		"status":422,
		"detail":"Failed Validation"
	}`
	apiError := &twizo.APIValidationError{}
	err := json.Unmarshal([]byte(cannedResponse), apiError)
	if err != nil {
		t.Fatal(err)
	}
	if apiError.Detail() != "Failed Validation" {
		t.Fatalf(
			"Invalid detail expecting [Failed Validation] got [%v]",
			apiError.Detail(),
		)
	}
}

func TestAPIErrorVerificationFieldResponse(t *testing.T) {
	cannedResponse := `{
    "errorCode":2,
    "type":"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
    "title":"Unprocessable Entity",
    "status":422,
    "detail":"AllowedTypes is empty after validation."
  }`
	apiError := &twizo.APIValidationError{}
	err := json.Unmarshal([]byte(cannedResponse), apiError)
	if err != nil {
		t.Fatal(err)
	}
	if apiError.Title() != "Unprocessable Entity" {
		t.Fatalf(
			"Invalid detail expecting [Unprocessable Entity] got [%v]",
			apiError.Title(),
		)
	}
	if apiError.Detail() != "AllowedTypes is empty after validation." {
		t.Fatalf(
			"Invalid detail expecting [AllowedTypes is empty after validation.] got [%v]",
			apiError.Detail(),
		)
	}
}
