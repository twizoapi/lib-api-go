package twizo_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func init() {
	twizo.APIKey = TestAPIKey
	twizo.RegionCurrent = TestRegion
}

func TestRegistrationWidgetSessionInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.RegistrationWidgetSessionResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestRegistrationWidgetSessionRequestExtra(t *testing.T) {
	sessionRequest := twizo.NewRegistrationWidgetSessionRequest()
	sessionRequest.SetBackupCodeIdentifier("backupCodeId")
	if "backupCodeId" != sessionRequest.GetBackupCodeIdentifier() {
		t.Fatalf(
			"Invalid BackupCodeID expecting [%s] got [%s]",
			"backupCodeId",
			sessionRequest.GetBackupCodeIdentifier(),
		)
	}
	sessionRequest.SetTotpIdentifier("totpCodeId")
	if "totpCodeId" != sessionRequest.GetTotpIdentifier() {
		t.Fatalf(
			"Invalid TotpCodeId expecting [%s] got [%s]",
			"totpCodeId",
			sessionRequest.GetTotpIdentifier(),
		)
	}
}

func TestRegistrationWidgetSessionRequest(t *testing.T) {
	data := struct {
		Issuer string
	}{
		Issuer: "issuer",
	}

	sessionRequest := twizo.NewRegistrationWidgetSessionRequest()
	sessionRequest.SetIssuer(data.Issuer)
	if data.Issuer != sessionRequest.GetIssuer() {
		t.Fatalf(
			"Invalid Issuer expecting [%s] got [%s]",
			data.Issuer,
			sessionRequest.GetIssuer(),
		)
	}

	// please not that Order is arrays must be maintaned as the 1st has perference over the 2nd
	sessionRequest.SetAllowedTypes([]string{"line", "sms"})
	if !reflect.DeepEqual(sessionRequest.GetAllowedTypes(), twizo.VerificationTypes{twizo.VerificationTypeLine, twizo.VerificationTypeSms}) {
		t.Fatalf(
			"Slices should be the same expected [%#v] got [%#v]",
			twizo.VerificationTypes{twizo.VerificationTypeSms, twizo.VerificationTypeLine},
			sessionRequest.GetAllowedTypes(),
		)
	}
	sessionRequest.SetAllowedTypes(twizo.VerificationTypes{twizo.VerificationTypeTelegram})
	if !reflect.DeepEqual(sessionRequest.GetAllowedTypes(), twizo.VerificationTypes{twizo.VerificationTypeTelegram}) {
		t.Fatalf(
			"Slices should be the same expected [%#v] got [%#v]",
			twizo.VerificationTypes{twizo.VerificationTypeTelegram},
			sessionRequest.GetAllowedTypes(),
		)
	}
}

func TestRegistrationWidgetSessionRequestSubmit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
		"sessionToken":"{{ .SessionToken }}",
		"applicationTag":"{{ .ApplicationTag }}",
		"requestedTypes":["sms"],
		"registeredTypes":["sms"],
		"allowedTypes":["sms"],
		"recipient":"{{ .Number }}",
		"totpIdentifier":null,
		"backupCodeIdentifier":null,
		"issuer":null,
		"language":"en",
		"status":"no status",
		"statusCode":0,
		"createdDateTime":"{{.CreatedDateTime}}",
		"_links":{
			"self":{
				"href":"https:\/\/{{ .Host }}\/v1\/widget-register-verification\/session\/{{ .SessionToken }}"
			}
		}
	}`
	data := struct {
		SessionToken    string
		ApplicationTag  string
		Number          string
		Host            string
		CreatedDateTime string
	}{
		SessionToken:    "test-01_11111110_wid5b45a482a62dc5.0457557354d4",
		ApplicationTag:  "app-tag",
		Number:          "6100000000",
		Host:            twizo.GetHostForRegion(twizo.RegionCurrent),
		CreatedDateTime: time.Now().UTC().Format(time.RFC3339),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	err = HTTPMockSend(
		http.MethodPost,
		fmt.Sprintf("https://%s/%s/widget-register-verification/session", data.Host, twizo.ClientAPIVersion),
		http.StatusCreated,
		b,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	session := twizo.NewRegistrationWidgetSessionRequest()
	session.SetRecipient(twizo.Recipient(data.Number))
	if session.GetRecipient() != twizo.Recipient(data.Number) {
		t.Fatalf("Invalid session recipient expecting [%v] got [%v]", twizo.Recipient(data.Number), session.GetRecipient())
	}

	response, err := session.Submit()
	if err != nil {
		t.Fatal(err)
	}

	// only check for now check sessionToken
	if response.GetSessionToken() != data.SessionToken {
		t.Fatalf("Invalid session token expecting [%s] got [%v]", data.SessionToken, response.GetSessionToken())
	}
	if response.GetStatusCode() != 0 {
		t.Fatalf("Invalid status code expecting [0] got [%v]", response.GetStatusCode())
	}
	if response.GetStatusMsg() != "no status" {
		t.Fatalf("Invalid status message expecting [no status] got [%v]", response.GetStatusMsg())
	}
	if response.GetLanguage() != "en" {
		t.Fatalf("Invalid language expecting [en] got [%v]", response.GetLanguage())
	}
	if response.GetBackupCodeIdentifier() != "" {
		t.Fatalf("Invalid backupcode expecting empty string got [%v]", response.GetBackupCodeIdentifier())
	}
	if response.GetApplicationTag() != data.ApplicationTag {
		t.Fatalf("Invalid applicationtag expecting [%s] got [%v]", data.ApplicationTag, response.GetApplicationTag())
	}
	if response.GetRecipient() != twizo.Recipient(data.Number) {
		t.Fatalf("Invalid recipient expecting [%v] got [%v]", twizo.Recipient(data.Number), response.GetRecipient())
	}
	if response.GetCreateDateTime().UTC().Format(time.RFC3339) != data.CreatedDateTime {
		t.Fatalf(
			"Invalid create time expecting [%v] got [%v]",
			data.CreatedDateTime,
			response.GetCreateDateTime().UTC().Format(time.RFC3339),
		)
	}
}
