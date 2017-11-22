package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"
	"testing"

	"encoding/json"
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
)

func init() {
	twizo.APIKey = TestApiKey
	twizo.RegionCurrent = TestRegion
}

func TestVerificationNew(t *testing.T) {
	verificationRequest := twizo.NewVerificationRequest(twizo.Recipient("0000000000"))

	_, err := json.Marshal(verificationRequest)
	if err != nil {
		t.Fatal(err)
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

	httpmock.RegisterResponder(http.MethodPost,
		fmt.Sprintf("https://%s/%s/verification/submit", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
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
	if response.GetMessageID() != "1" {
		t.Fatalf("Invalid message id expected [1] got [%v]", response.GetMessageID())
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
		fmt.Sprintf("https://%s/%s/verification/submit", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusCreated, TestResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	request := twizo.NewVerificationRequest(twizo.Recipient("0000000000"))

	response, err := request.Submit()
	if err != nil {
		t.Fatal(err)
	}
	if response.GetMessageID() != "1" {
		t.Fatalf("Invalid message id expected [1] got [%v]", response.GetMessageID())
	}
}
