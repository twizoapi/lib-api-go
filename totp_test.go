package twizo_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"

	"github.com/bitly/go-simplejson"
	"gopkg.in/jarcoal/httpmock.v1"
)

func init() {
	twizo.APIKey = TestAPIKey
	twizo.RegionCurrent = TestRegion
}

func TestTotpInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.TotpResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestTotpInvalidResponseUrl(t *testing.T) {
	// our basic response we will need to add more here
	const tpl = `{
		"identifier": "identifier",
		"issuer": "issuer",
		"uri": "%",
		"verification": null,
		"_links": {
			"self": {
				"href": "https://{{ .Host }}/v1/totp/tester"
			}
		}
	}`

	data := struct {
		Host string
	}{
		Host: twizo.GetHostForRegion(twizo.RegionCurrent),
	}
	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	response := &twizo.TotpResponse{}
	err = response.UnmarshalJSON(b)
	// error must be url.Error
	if _, ok := err.(*url.Error); !ok {
		t.Fatalf(
			"Expecting url.error for invalid url got [%#v]",
			err,
		)
	}
}

func TestTotpRequest(t *testing.T) {
	data := struct {
		Identifier string
		Issuer     string
	}{
		Identifier: "identifier",
		Issuer:     "issuer",
	}

	request := twizo.NewTotpRequest(data.Identifier)
	request.SetIssuer(data.Issuer)

	if request.GetIdentifier() != data.Identifier {
		t.Fatalf(
			"Invalid identifier expecting [%#v] got [%#v]",
			data.Identifier,
			request.GetIdentifier(),
		)
	}
	if request.GetIssuer() != data.Issuer {
		t.Fatalf(
			"Invalid issuer expecting [%#v] got [%#v]",
			data.Issuer,
			request.GetIssuer(),
		)
	}
}

func TestTotpDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	data := struct {
		Identifier string
		Host       string
	}{
		Identifier: "test",
		Host:       twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	err := HTTPMockSend(
		http.MethodDelete,
		fmt.Sprintf(
			"https://%s/%s/totp/%s",
			data.Host,
			twizo.ClientAPIVersion,
			data.Identifier,
		),
		http.StatusNoContent,
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = twizo.TotpDelete(data.Identifier)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTotpCreateResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
		"identifier": "{{ .Identifier }}",
		"issuer": "{{ .Issuer }}",
		"uri": "{{ .URI }}",
		"verification": null,
		"_links": {
			"self": {
				"href": "https://{{ .Host }}/v1/totp/{{ .Identifier }}"
			}
		}
	}`

	data := struct {
		Identifier string
		Issuer     string
		Token      string
		Secret     string
		URI        *url.URL
		Host       string
	}{
		Identifier: "test",
		Issuer:     "test",
		Token:      "test",
		Secret:     "TESTKEYTESTKEYTESTKEYTESTKEYTESTKEYTESTKEYTESTKEYTES",
		Host:       twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	data.URI, _ = url.Parse("otpauth://totp")
	data.URI.Path = fmt.Sprintf("/%s:%s", data.Issuer, data.Identifier)
	q := data.URI.Query()
	q.Set("issuer", data.Issuer)
	q.Set("secret", data.Secret)
	data.URI.RawQuery = q.Encode()

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	err = HTTPMockSend(
		http.MethodPost,
		fmt.Sprintf(
			"https://%s/%s/totp",
			data.Host,
			twizo.ClientAPIVersion,
		),
		http.StatusCreated,
		b,
		func(req *http.Request) error {

			js, funcErr := simplejson.NewFromReader(req.Body)
			if funcErr != nil {
				return funcErr
			}
			// must contain 2 elements
			v, funcErr := js.Map()
			if funcErr != nil {
				return funcErr
			}
			if len(v) != 2 {
				return fmt.Errorf("expected 2 elements got [%d]", len(v))
			}

			// check for identifier
			if identifier, ok := js.CheckGet("identifier"); ok {
				if identifier.MustString() != data.Identifier {
					return fmt.Errorf("identifier expected to be [%s] got [%v]", data.Identifier, identifier.MustString())
				}
			} else {
				return fmt.Errorf("identifier not defined")
			}

			// check for issuer
			if issuer, ok := js.CheckGet("issuer"); ok {
				if issuer.MustString() != data.Issuer {
					return fmt.Errorf("issuer expected to be [%s] got [%v]", data.Issuer, issuer.MustString())
				}
			} else {
				return fmt.Errorf("issuer not defined")
			}

			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.TotpCreate(data.Identifier, data.Issuer)
	if err != nil {
		t.Fatal(err)
	}

	// there is no verification response but it should not panic
	if response.GetVerificationResponse().IsTokenSuccess() {
		t.Fatal("verification should not be a succes")
	}

	if !reflect.DeepEqual(response.GetURL(), data.URI) {
		t.Fatalf(
			"Invalid uri expecting [%#v] got [%#v]",
			data.URI,
			response.GetURL(),
		)
	}
	if *response.GetURLSecret() != data.Secret {
		t.Fatalf(
			"Invalid Secret expecting [%#v] got [%#v]",
			data.Secret,
			response.GetURLSecret(),
		)
	}
}

func TestTotpCheckResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
		"identifier": "{{ .Identifier }}",
		"issuer": "{{ .Issuer }}",
		"uri": null,
		"verification": null,
		"_links": {
			"self": {
				"href": "https://{{ .Host }}/v1/totp/{{ .Identifier }}"
			}
		}
	}`

	data := struct {
		Identifier string
		Issuer     string
		Host       string
	}{
		Identifier: "identifier",
		Issuer:     "issuer",
		Host:       twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}
	_ = HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf(
			"https://%s/%s/totp/%s",
			data.Host,
			twizo.ClientAPIVersion,
			data.Identifier,
		),
		http.StatusOK,
		b,
		nil,
	)

	response, err := twizo.TotpCheck(data.Identifier)
	if err != nil {
		t.Fatal(err)
	}

	if response.GetIdentifier() != data.Identifier {
		t.Fatalf(
			"Invalid identifier expecting [%#v] got [%#v]",
			data.Identifier,
			response.GetIdentifier(),
		)
	}
	if response.GetIssuer() != data.Issuer {
		t.Fatalf(
			"Invalid issuer expecting [%#v] got [%#v]",
			data.Issuer,
			response.GetIssuer(),
		)
	}
	if response.GetVerificationResponse().IsTokenSuccess() {
		t.Fatal("verification should not be a succes")
	}
	if response.GetURL() != nil {
		t.Fatal("url should be nil, as none was sent in the response")
	}
	if response.GetURLSecret() != nil {
		t.Fatal("there should be no secret in the url")
	}
}

func TestTotpVerifyResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
		"identifier": "{{ .Identifier }}",
		"issuer": "{{ .Issuer }}",
		"uri": null,
		"_embedded": {
			"verification": {
				"applicationTag": "Test",
				"bodyTemplate": null,
				"createdDateTime": "2018-02-01T08:32:22+00:00",
				"dcs": null,
				"issuer": "",
				"language": null,
				"messageId": "eu-01-2.24909.ver5a72d096d55c27.50239050",
				"reasonCode": null,
				"recipient": "",
				"salesPrice": null,
				"salesPriceCurrencyCode": null,
				"sender": null,
				"senderNpi": null,
				"senderTon": null,
				"sessionId": "",
				"status": "success",
				"statusCode": 1,
				"tag": null,
				"tokenLength": null,
				"tokenType": null,
				"type": "totp",
				"validity": null,
				"validUntilDateTime": null,
				"voiceSentence": null,
				"webHook": null,
				"_links": {
					"self": {
						"href": "https://{{ .Host }}/v1/verification/submit/eu-01-2.24909.ver5a72d096d55c27.50239050"
					}
				}
			}
		},
		"_links": {
			"self": {
				"href": "https://{{ .Host }}/v1/totp/{{ .Identifier }}"
			}
		}
	}`

	data := struct {
		Identifier string
		Issuer     string
		Token      string
		Host       string
	}{
		Identifier: "test",
		Issuer:     "test",
		Token:      "test",
		Host:       twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}
	_ = HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf(
			"https://%s/%s/totp/%s?token=%s",
			data.Host,
			twizo.ClientAPIVersion,
			data.Identifier,
			data.Token,
		),
		http.StatusOK,
		b,
		nil,
	)

	response, err := twizo.TotpVerify(data.Identifier, data.Token)
	if err != nil {
		t.Fatal(err)
	}

	if !response.GetVerificationResponse().IsTokenSuccess() {
		t.Fatal("verification should be a succes")
	}
	if response.GetIssuer() != data.Issuer {
		t.Fatalf(
			"Invalid issuer expecting [%#v] got [%#v]",
			data.Issuer,
			response.GetIssuer(),
		)
	}
}
