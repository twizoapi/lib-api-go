package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

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

func TestBackupCodeInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.BackupCodeResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestBackupCodeCreateAlreadyExists(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cannedResponse := twizo.NewAPIError("Conflict", http.StatusConflict)

	err := HTTPMockSend(
		http.MethodPost,
		fmt.Sprintf("https://%s/%s/backupcode", twizo.GetHostForRegion(twizo.RegionCurrent), twizo.ClientAPIVersion),
		http.StatusConflict,
		cannedResponse,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.BackupCodeCreate("example@example.com")
	if err != nil {
		t.Fatal(err)
	}

	if !response.AlreadyExists() {
		t.Fatalf(
			"Invalid alreadyExists flag expecting [%v] got [%v]",
			true,
			response.AlreadyExists(),
		)
	}
}

func TestBackupCodeUpdate(t *testing.T) {
	doTestBackupCodeCreate(t, true)
}

func TestBackupCodeCreate(t *testing.T) {
	doTestBackupCodeCreate(t, false)
}

func doTestBackupCodeCreate(t *testing.T, update bool) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	tpl := `{
    "identifier": "{{ .Identifier }}",
    "amountOfCodesLeft": {{ .AmountOfCodesLeft }},
    "codes": {{ .Codes }},
    "createdDateTime": "{{ .CreatedDateTime }}",
    "_links": {
        "self": {
            "href": "https://{{ .Host }}/v1/backupcode/{{ .Identifier }}"
        }
    }
  }`

	codes := []string{
		"01234560",
		"01234561",
		"01234562",
		"01234563",
		"01234564",
		"01234565",
		"01234566",
		"01234567",
		"01234568",
		"01234569",
	}
	jsonCodes, _ := json.Marshal(codes)
	data := struct {
		Identifier        string
		AmountOfCodesLeft int
		Codes             string
		CreatedDateTime   string
		Host              string
	}{
		Identifier:        "example@example.com",
		AmountOfCodesLeft: len(codes),
		Codes:             string(jsonCodes),
		CreatedDateTime:   time.Now().UTC().Format(time.RFC3339),
		Host:              twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	url := "backupcode"
	method := http.MethodPost
	expect := http.StatusCreated
	if update {
		url = fmt.Sprintf("backupcode/%s", data.Identifier)
		method = http.MethodPut
		expect = http.StatusOK
	}

	err = HTTPMockSend(
		method,
		fmt.Sprintf("https://%s/%s/%s", data.Host, twizo.ClientAPIVersion, url),
		expect,
		b,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	var response = &twizo.BackupCodeResponse{}
	if update {
		response, err = twizo.BackupCodeUpdate(data.Identifier)
	} else {
		response, err = twizo.BackupCodeCreate(data.Identifier)
	}

	if err != nil {
		t.Fatal(err)
	}

	if response.GetIdentifier() != data.Identifier {
		t.Fatalf(
			"Invalid identifier expecting [%s] got [%v]",
			data.Identifier,
			response.GetIdentifier(),
		)
	}

	if response.GetAmountOfCodesLeft() != data.AmountOfCodesLeft {
		t.Fatalf(
			"Invalid amount of codes left expecting [%d] got [%v]",
			data.AmountOfCodesLeft,
			response.GetAmountOfCodesLeft(),
		)
	}

	if response.GetCreateDateTime().UTC().Format(time.RFC3339) != data.CreatedDateTime {
		t.Fatalf(
			"Invalid create time expecting [%v] got [%v]",
			data.CreatedDateTime,
			response.GetCreateDateTime().UTC().Format(time.RFC3339),
		)
	}

	if !reflect.DeepEqual(response.GetCodes(), codes) {
		t.Fatalf(
			"Invalid codes shoud be the same expecting [%v] got [%v]",
			codes,
			response.GetCodes(),
		)
	}

	if response.AlreadyExists() {
		t.Fatalf(
			"Invalid alreadyExists flag expecting [%v] got [%v]",
			false,
			response.AlreadyExists(),
		)
	}
}

func TestBackupCodeDelete(t *testing.T) {
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
			"https://%s/%s/backupcode/%s",
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

	err = twizo.BackupCodeDelete(data.Identifier)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackupCodeVerify(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tpl := `{
		"identifier":"{{ .Identifier }}",
		"amountOfCodesLeft": {{.AmountOfCodesLeft}},
		"codes":[],
		"createdDateTime":"2018-02-14T09:32:26+00:00",
		"_embedded":{
			"verification":{
				"applicationTag":"applicationTag",
				"bodyTemplate":null,
				"createdDateTime":"2018-02-14T09:32:48+00:00",
				"dcs":null,
				"issuer":"",
				"language":null,
				"messageId":"eu-01-2.10570.ver5a8402404e8b12.86770577",
				"reasonCode":null,
				"recipient":"",
				"salesPrice":null,
				"salesPriceCurrencyCode":null,
				"sender":null,
				"senderNpi":null,
				"senderTon":null,
				"sessionId":"",
				"status":"success",
				"statusCode":1,
				"tag":null,
				"tokenLength":null,
				"tokenType":null,
				"type":"backupcode",
				"validity":null,
				"validUntilDateTime":null,
				"voiceSentence":null,
				"webHook":null,
				"_links":{
					"self":{
						"href":"https:\/\/api-eu-01.twizo.com\/v1\/verification\/submit\/eu-01-2.10570.ver5a8402404e8b12.86770577"
					}
				}
			}
		},
		"_links":{
			"self":{
				"href":"https:\/\/api-eu-01.twizo.com\/v1\/backupcode\/12345"
			}
		}
	}`

	data := struct {
		Identifier        string
		Token             string
		Host              string
		AmountOfCodesLeft int
	}{
		Identifier:        "identifier",
		Token:             "12345",
		Host:              twizo.GetHostForRegion(twizo.RegionCurrent),
		AmountOfCodesLeft: 9,
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf(
			"https://%s/%s/backupcode/%s?token=%s",
			data.Host,
			twizo.ClientAPIVersion,
			data.Identifier,
			url.QueryEscape(data.Token),
		),
		http.StatusOK,
		b,
		nil,
	)

	response, err := twizo.BackupCodeVerify(data.Identifier, data.Token)
	if err != nil {
		t.Fatal(err)
	}
	if response.GetAmountOfCodesLeft() != data.AmountOfCodesLeft {
		t.Fatalf(
			"Invalid amount of codes left expecting [%d] got [%v]",
			data.AmountOfCodesLeft,
			response.GetAmountOfCodesLeft(),
		)
	}

	// test some error states
	HTTPMockSend(
		http.MethodGet,
		fmt.Sprintf(
			"https://%s/%s/backupcode/%s?token=%s",
			data.Host,
			twizo.ClientAPIVersion,
			data.Identifier,
			url.QueryEscape(data.Token),
		),
		http.StatusLocked,
		b,
		nil,
	)

	response, err = twizo.BackupCodeVerify(data.Identifier, data.Token)
	if err == nil {
		t.Fatalf("Expecting error here")
	}
	_, ok := err.(*twizo.ClientError)
	if ok == false {
		t.Fatalf("Receiving invalid json expecting [twizo.ClientError] got [%#v]", err)
		return
	}
}
