package twizo_test

import (
	"testing"

	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"

	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNewUrl(t *testing.T) {
	apiURL, err := twizo.GetURLFor("test")
	if err != nil {
		t.Fatal(err)
	}

	if apiURL.Scheme != "https" {
		t.Fatalf("Incorrect scheme expecting [https] got [%v]", apiURL.Scheme)
	}

	if apiURL.Path != fmt.Sprintf("/%s/test", twizo.ClientAPIVersion) {
		t.Fatalf("Incorrect path expecting [/%s/test] got [%v]", twizo.ClientAPIVersion, apiURL.Path)
	}
}

func TestNewRequest(t *testing.T) {
	apiURL, err := twizo.GetURLFor("test")
	if err != nil {
		t.Fatal(err)
	}

	req, err := twizo.GetClient(TestRegion, TestAPIKey).NewRequest("GET", apiURL, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "GET" {
		t.Fatalf("Incorrect Authorization header expecting [GET] got [%v]", req.Method)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", twizo.ClientAuthUser, TestAPIKey)))
	if req.Header.Get("Authorization") != fmt.Sprintf("Basic %s", auth) {
		t.Fatalf("Incorrect Authorization header expecting [Basic %s] got [%v]", auth, req.Header.Get("Authorization"))
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Fatalf("Incorrect Accept header expecting [application/json] got [%v]", req.Header.Get("Accept"))
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Incorrect Content-Type header expecting [application/json] got [%v]", req.Header.Get("Accept"))
	}

	if req.URL.Host != twizo.GetHostForRegion(TestRegion) {
		t.Fatalf("Incorrect URL expecting [%v] got [%v]", twizo.GetHostForRegion(TestRegion), req.URL.Host)
	}
}

func TestReceivingBadJsonRequests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//httpmock.RegisterResponder("POST", fmt.Sprintf("https://%s/v1/broken", twizo.GetHostForRegion(twizo.RegionCurrent)),
	//	httpmock.NewStringResponder(200, `{"messageId": 1}`))

	TestRequest := struct {
		Valid string `json:"valid"`
	}{
		Valid: "1",
	}

	TestResponse := struct {
		Valid string `json:"valid"`
	}{}

	client := twizo.GetClient(TestRegion, TestAPIKey)

	// Test Returning bad json
	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://%s/%s/broken", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
		httpmock.NewStringResponder(http.StatusConflict, `this is not json`),
	)

	apiURL, err := twizo.GetURLFor("/broken")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Call("POST", apiURL, TestRequest, http.StatusConflict, TestResponse)
	if err == nil {
		t.Fatal("Recieving invalid json expecting [error] got [nil]")
	}
	_, ok := err.(*json.SyntaxError)
	if ok == false {
		t.Fatalf("Receiving invalid json expecting [json.SyntaxError] got [%#v]", err)
		return
	}
}

func TestReceivingUnexpectedCodeRequests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := twizo.GetClient(TestRegion, TestAPIKey)

	// Test Returning bad json
	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://%s/%s/broken", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
		httpmock.NewStringResponder(http.StatusConflict, `this is not json`),
	)

	apiURL, err := twizo.GetURLFor("/broken")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Call("POST", apiURL, nil, http.StatusOK, nil)
	if err == nil {
		t.Fatal("Recieving unexpected status code, expecting [error] got [nil]")
	}
	_, ok := err.(*twizo.ClientError)
	if ok == false {
		t.Fatalf("Receiving invalid json expecting [twizo.ClientError] got [%#v]", err)
		return
	}
}
