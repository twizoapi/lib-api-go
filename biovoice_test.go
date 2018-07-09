package twizo_test

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

func TestBioVoiceInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.BioVoiceResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestBioVoiceCheckRegistration(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
		"createdDateTime": "2018-02-02T15:04:15+00:00",
		"language": null,
		"reasonCode": null,
		"recipient": "{{ .Recipient }}",
		"registrationId": "{{ .RegistrationID }}",
		"salesPrice": null,
		"salesPriceCurrencyCode": null,
		"status": "no status",
		"statusCode": 0,
		"voiceSentence": "Verify me with my voicepin",
		"webHook": null,
		"_links": {
			"self": {
				"href": "https://{{ .Host }}/v1/biovoice/registration/{{ .RegistrationID }}"
			}
		}
	}`

	data := struct {
		Recipient      twizo.Recipient
		RegistrationID string
		Host           string
	}{
		Recipient:      twizo.Recipient("0000000000"),
		RegistrationID: "00000.B000fff000fff0000.00000000",
		Host:           twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}
	err = HTTPMockSend(
		http.MethodPost,
		fmt.Sprintf(
			"https://%s/%s/biovoice/registration",
			data.Host,
			twizo.ClientAPIVersion,
		),
		http.StatusCreated,
		b,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.BioVoiceCreateRegistration(data.Recipient)
	if err != nil {
		t.Fatal(err)
	}

	if response.GetRecipient() != data.Recipient {
		t.Fatalf(
			"Invalid GetRecipient expecting [%v] got [%v]",
			data.Recipient,
			response.GetRecipient(),
		)

	}
}
