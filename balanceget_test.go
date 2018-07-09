package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	"encoding/json"
	"errors"
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

func TestBalanceInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.BalanceGetResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestBalanceInvalidJson(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	err := HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf("https://%s/%s/wallet/getbalance", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
		http.StatusOK,
		"This is invalid json",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.BalanceGet()
	if response != nil {
		t.Fatal(errors.New("testbalanceinvalidjson: not expecting result on error"))
	}
	if err == nil {
		t.Fatal(errors.New("testbalanceinvalidjson: expecting error on invalid json"))
	}
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestBalanceGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
    "credit": {{ .Credit }},
    "currencyCode":"{{ .CurrencyCode }}",
    "freeVerifications":{{.FreeVerifications}},
    "wallet":"{{ .Wallet }}",
	  "alarmLimit":"{{.AlarmLimit}}"
  }`

	data := struct {
		Credit            float32
		CurrencyCode      string
		Host              string
		FreeVerifications int
		Wallet            string
		AlarmLimit        string
	}{
		Credit:            1.666,
		CurrencyCode:      "eur",
		FreeVerifications: 1,
		Wallet:            "My Test Wallet",
		AlarmLimit:        "0.50",
		Host:              twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	err = HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf("https://%s/%s/wallet/getbalance", data.Host, twizo.ClientAPIVersion),
		http.StatusOK,
		b,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.BalanceGet()
	if err != nil {
		t.Fatal(err)
	}

	if response.GetCurrencyCode() != data.CurrencyCode {
		t.Fatalf(
			"Invalid currency code expecting [%s] got [%v]",
			data.CurrencyCode,
			response.GetCurrencyCode(),
		)
	}
	if *response.GetAlarmLimit() != data.AlarmLimit {
		t.Fatalf(
			"Invalid alarm limit expecting [%s] got [%v]",
			data.AlarmLimit,
			response.GetAlarmLimit(),
		)
	}
	if response.GetCredit() != data.Credit {
		t.Fatalf(
			"Invalid credit expecting [%f] got [%v]",
			data.Credit,
			response.GetCredit(),
		)
	}
	if response.GetWallet() != data.Wallet {
		t.Fatalf(
			"Invalid wallet expecting [%s] got [%v]",
			data.Wallet,
			response.GetWallet(),
		)
	}
	if response.GetFreeVerifications() != data.FreeVerifications {
		t.Fatalf(
			"Invalid free verifications expecting [%v] got [%v]",
			data.FreeVerifications,
			response.GetFreeVerifications(),
		)
	}
}
