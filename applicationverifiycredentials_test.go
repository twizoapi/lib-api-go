package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	"encoding/json"
	"errors"
	"testing"

	// twizo "github.com/twizoapi/lib-api-go"
	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"

	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/jarcoal/httpmock.v1"
)

func init() {
	twizo.APIKey = TestAPIKey
	twizo.RegionCurrent = TestRegion
}

func TestApplicationVerifyInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.ApplicationVerifyCredentialsResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestInvalidUrl(t *testing.T) {
	defer func() { twizo.GetURLFor = twizo.GetURLForOriginal }()
	twizo.GetURLFor = func(string) (*url.URL, error) { return nil, errors.New("fail") }
	response, err := twizo.ApplicationVerifyCredentials()
	if err == nil {
		t.Fatalf(
			"Expecting error got [%#v]",
			err,
		)
	}
	if response != nil {
		t.Fatalf(
			"Expecting null response got [%#v]",
			response,
		)
	}
}

func TestInvalidKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cannedResponse := twizo.NewAPIError("Unauthorized", http.StatusUnauthorized)

	err := HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf(
			"https://%s/%s/application/verifycredentials",
			twizo.GetHostForRegion(twizo.RegionCurrent),
			twizo.ClientAPIVersion,
		),
		http.StatusUnauthorized,
		cannedResponse,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.ApplicationVerifyCredentials()
	if err != nil {
		t.Fatal(err)
	}

	if response.IsKeyValid() {
		t.Fatalf(
			"Invalid test key status [%v] got [%v]",
			true,
			response.IsTestKey(),
		)
	}
}

func TestValidKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const tpl = `{"applicationTag":"{{ .ApplicationTag }}","isTestKey":{{ .IsTestKey }}}`

	data := struct {
		ApplicationTag string
		IsTestKey      bool
		Host           string
	}{
		ApplicationTag: "ApplicationTag",
		IsTestKey:      true,
		Host:           twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	err = HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf("https://%s/%s/application/verifycredentials", data.Host, twizo.ClientAPIVersion),
		http.StatusOK,
		b,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.ApplicationVerifyCredentials()
	if err != nil {
		t.Fatal(err)
	}

	if response.IsTestKey() != data.IsTestKey {
		t.Fatalf(
			"Invalid test key status [%v] got [%v]",
			data.IsTestKey,
			response.IsTestKey(),
		)
	}
	if response.GetApplicationTag() != data.ApplicationTag {
		t.Fatalf(
			"Invalid free verifications expecting [%s] got [%v]",
			data.ApplicationTag,
			response.GetApplicationTag(),
		)
	}
	if response.IsKeyValid() != true {
		t.Fatalf(
			"Invalid key should be valid [%s] got [%v]",
			data.ApplicationTag,
			response.GetApplicationTag(),
		)
	}
}
