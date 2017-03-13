package twizo_test

import (
	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"
	"testing"

	"bytes"
	"encoding/base64"
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
)

func TestNewUrl(t *testing.T) {
	apiURL, err := twizo.GetURLFor("test")
	if err != nil {
		t.Fatal(err)
	}

	if apiURL.Scheme != "https" {
		t.Fatalf("Incorrect scheme expected [https] got [%v]", apiURL.Scheme)
	}

	if apiURL.Path != fmt.Sprintf("/%s/test", twizo.ClientAPIVersion) {
		t.Fatalf("Incorrect path expected [/%s/test] got [%v]", twizo.ClientAPIVersion, apiURL.Path)
	}
}

func TestNewRequest(t *testing.T) {

	apiURL, err := twizo.GetURLFor("test")
	if err != nil {
		t.Fatal(err)
	}

	req, err := twizo.GetClient(TestRegion, TestApiKey).NewRequest("GET", apiURL, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "GET" {
		t.Fatalf("Incorrect Authorization header expected [GET] got [%v]", req.Method)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", twizo.ClientAuthUser, TestApiKey)))
	if req.Header.Get("Authorization") != fmt.Sprintf("Basic %s", auth) {
		t.Fatalf("Incorrect Authorization header expected [Basic %s] got [%v]", auth, req.Header.Get("Authorization"))
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Fatalf("Incorrect Accept header expected [application/json] got [%v]", req.Header.Get("Accept"))
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Incorrect Content-Type header expected [application/json] got [%v]", req.Header.Get("Accept"))
	}

	if req.URL.Host != twizo.GetHostForRegion(TestRegion) {
		t.Fatalf("Incorrect URL expected [%v] got [%v]", twizo.GetHostForRegion(TestRegion), req.URL.Host)
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

	client := twizo.GetClient(TestRegion, TestApiKey)

	// Test Returning bad json
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://%s/%s/broken", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
		httpmock.NewStringResponder(http.StatusConflict, `this is not json`))

	apiURL, err := twizo.GetURLFor("/broken")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Call("POST", apiURL, TestRequest, http.StatusConflict, TestResponse)
	if err == nil {
		t.Fatal("Recieving invalid json expects [error] got [nil]")
	}
	_, ok := err.(*twizo.JSONClientError)
	if ok == false {
		t.Fatalf("Receiving invalid json expects [twizo.JsonClientError] got [%#v]", err)
		return
	}
}

func TestReceivingUnexpectedCodeRequests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := twizo.GetClient(TestRegion, TestApiKey)

	// Test Returning bad json
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://%s/%s/broken", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
		httpmock.NewStringResponder(http.StatusConflict, `this is not json`))

	apiURL, err := twizo.GetURLFor("/broken")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Call("POST", apiURL, nil, http.StatusOK, nil)
	if err == nil {
		t.Fatal("Recieving unexpected status code expects [error] got [nil]")
	}
	_, ok := err.(*twizo.ClientError)
	if ok == false {
		t.Fatalf("Receiving invalid json expects [twizo.ClientError] got [%#v]", err)
		return
	}
}
