package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func init() {
	twizo.APIKey = TestAPIKey
	twizo.RegionCurrent = TestRegion
}

func TestNumberLookupInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.NumberLookupResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%v]",
			err,
		)
	}
}

func TestNumberLookupSubmit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
    "_links": {
        "self": {
            "href": "https://{{ .Host }}/v1/numberlookup/submit"
        }
    },
    "_embedded": {
        "items": [
            {
                "applicationTag": "{{.ApplicationTag}}",
                "callbackUrl": null,
                "countryCode": null,
                "createdDateTime": "2017-03-13T14:37:35+00:00",
                "imsi": null,
                "isPorted": "Unknown",
                "isRoaming": "Unknown",
                "messageId": "{{ .MessageID }}",
                "msc": null,
                "networkCode": null,
                "number": "{{ .Number }}",
                "operator": null,
                "reasonCode": null,
                "resultTimestamp": null,
                "resultType": 0,
                "salesPrice": null,
                "salesPriceCurrencyCode": null,
                "status": "no status",
                "statusCode": 0,
                "tag": null,
                "validity": 259200,
                "validUntilDateTime": "2017-03-16T14:37:35+00:00",
                "_links": {
                    "self": {
                        "href": "https://{{ .Host }}/v1/numberlookup/submit/{{ .MessageID }}"
                    }
                }
            }
        ]
    },
    "total_items": 1
}`

	data := struct {
		MessageID      string
		Number         string
		ApplicationTag string
		Host           string
	}{
		MessageID:      "test-1.10314.sms58c16b15c261a5.18930279",
		Number:         "6100000000",
		ApplicationTag: "UnitTest",
		Host:           twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	err = HTTPMockSend(
		http.MethodPost,
		fmt.Sprintf(
			"https://%s/%s/numberlookup/submit",
			data.Host,
			twizo.ClientAPIVersion,
		),
		http.StatusCreated,
		b,
		func(req *http.Request) error {
			var arbitraryJSON map[string]interface{}
			receivedJSON, subErr := ioutil.ReadAll(req.Body)
			if subErr != nil {
				panic(subErr)
			}
			subErr = json.Unmarshal([]byte(receivedJSON), &arbitraryJSON)
			if err != nil {
				panic(subErr)
			}
			for key, value := range arbitraryJSON {
				switch key {
				case "numbers":
					values := value.([]interface{})
					if len(values) != 1 {
						return fmt.Errorf("too many elements in request for numbers expecting 1 got [%d]", len(values))
					}
					if values[0].(string) != data.Number {
						return fmt.Errorf("number element should be [%s] got [%v]", data.Number, values[0])
					}
				default:
					return fmt.Errorf("unexpected element [%s] -> [%v] in request", key, value)
				}
			}
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	response, err := twizo.NumberLookupSubmit([]twizo.Recipient{twizo.Recipient(data.Number)})
	if err != nil {
		t.Fatal(err)
	}
	items := response.GetItems()
	if len(items) != 1 {
		t.Fatalf("Expecting 1 returned element got [%d]", len(items))
	}

	// only check for now check messageId
	if items[0].GetMessageID() != data.MessageID {
		t.Fatalf("Invalid message id expecting [%s] got [%v]", data.MessageID, items[0].GetMessageID())
	}
	if items[0].GetNumber() != data.Number {
		t.Fatalf("Invalid message id expecting [%s] got [%v]", data.Number, items[0].GetNumber())
	}
	if items[0].GetApplicationTag() != data.ApplicationTag {
		t.Fatalf("Invalid application tag expecting [%s] got [%v]", data.ApplicationTag, items[0].GetApplicationTag())
	}
}
