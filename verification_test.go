package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	"encoding/json"
	"testing"

	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"

	"fmt"
	"net/http"

	"gopkg.in/jarcoal/httpmock.v1"
)

func init() {
	twizo.APIKey = TestAPIKey
	twizo.RegionCurrent = TestRegion
}

func TestVerificationInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.VerificationResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestVerificationNew(t *testing.T) {
	// from string
	_, err := twizo.NewVerificationRequest("0000000000")
	if err != nil {
		t.Fatal(err)
	}
	// from recipient
	_, err = twizo.NewVerificationRequest(twizo.Recipient("0000000000"))
	if err != nil {
		t.Fatal(err)
	}
	// from invalid => error
	_, err = twizo.NewVerificationRequest(1)
	if err == nil {
		t.Fatalf("expecting nil got [%#v]", err)
	}
	// from multiple (can only handle one recipient) => error
	_, err = twizo.NewVerificationRequest([]string{"one", "two"})
	if err == nil {
		t.Fatalf("expecting error got [%#v]", err)
	}
}

func TestVerificationSubmit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	TestResponse := struct {
		MessageID string `json:"messageId"`
	}{
		MessageID: "1",
	}

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf(
			"https://%s/%s/verification/submit",
			twizo.GetHostForRegion(twizo.RegionCurrent),
			twizo.ClientAPIVersion,
		),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusCreated, TestResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	response, err := twizo.VerificationSubmit(twizo.Recipient("0000000000"))
	if err != nil {
		t.Fatal(err)
	}
	if response.GetMessageID() != TestResponse.MessageID {
		t.Fatalf("Invalid message id expecting [%s] got [%v]", TestResponse.MessageID, response.GetMessageID())
	}
}

func TestVerificationSubmitAdvanced(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	TestResponse := struct {
		MessageID string `json:"messageId"`
	}{
		MessageID: "1",
	}

	httpmock.RegisterResponder(http.MethodPost,
		fmt.Sprintf(
			"https://%s/%s/verification/submit",
			twizo.GetHostForRegion(twizo.RegionCurrent),
			twizo.ClientAPIVersion,
		),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusCreated, TestResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	request, _ := twizo.NewVerificationRequest(twizo.Recipient("0000000000"))

	response, err := request.Submit()
	if err != nil {
		t.Fatal(err)
	}
	if response.GetMessageID() != TestResponse.MessageID {
		t.Fatalf("Invalid message id expecting [%s] got [%v]", TestResponse.MessageID, response.GetMessageID())
	}
}

/**
 * Tests for VerificationTypes
 */

func TestVerificationTypesFetchArray(t *testing.T) {
	cannedResponse := `["sms","push","totpCreate","biovoice","call","telegram","line","backupcode"]`
	testVerificationTypesFetch(t, cannedResponse)
}

func TestVerificationTypesFetchMap(t *testing.T) {
	// due to an api bug this can happen ;(, it's not within spec but we
	// support it for now an issue has been raised with Twizo
	cannedResponse := `{
		"0":"sms",
		"1":"push",
		"2":"totpCreate",
		"3":"biovoice",
		"4":"call",
		"5":"telegram",
		"6":"line",
		"7":"backupcode"
	}`
	testVerificationTypesFetch(t, cannedResponse)
}

func testVerificationTypesFetch(t *testing.T, response interface{}) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	err := HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf(
			"https://%s/%s/application/verification_types",
			twizo.GetHostForRegion(twizo.RegionCurrent),
			twizo.ClientAPIVersion,
		),
		http.StatusOK,
		response,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	verificationTypes := twizo.NewVerificationTypes()
	err = verificationTypes.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	if !verificationTypes.Has(twizo.VerificationType("sms")) {
		t.Fatal("invalid element not present expecting [sms] to be present")
	}
}

func TestVerificationTypesDuplicates(t *testing.T) {
	verificationTypes := twizo.NewVerificationTypes()
	verificationTypes.Add("sms")
	verificationTypes.Add("sms")
	json, err := json.Marshal(verificationTypes)
	if err != nil {
		t.Fatal(err)
	}
	if string(json) != `["sms"]` {
		t.Fatal("invalid duplicate elements not present expecting [\"sms\"] got", string(json))
	}
}
