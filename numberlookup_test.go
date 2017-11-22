package twizo_test

// go get gopkg.in/jarcoal/httpmock.v1

import (
	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"
	"testing"

	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
)

func init() {
	twizo.APIKey = TestApiKey
	twizo.RegionCurrent = TestRegion
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
                "applicationTag": "UnitTest",
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
		MessageID string
		Number    string
		Host      string
	}{
		MessageID: "test-1.10314.sms58c16b15c261a5.18930279",
		Number:    "6100000000",
		Host:      twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	_ = HttpMockSendJsonPostTo(
		fmt.Sprintf("https://%s/%s/numberlookup/submit", data.Host, twizo.ClientAPIVersion),
		http.StatusCreated,
		b,
	)

	response, err := twizo.NumberLookupSubmit([]twizo.Recipient{twizo.Recipient(data.Number)})
	if err != nil {
		t.Fatal(err)
	}
	items := response.GetItems()

	// only check for now check messageId
	if items[0].GetMessageID() != data.MessageID {
		t.Fatalf("Invalid message id expected [%s] got [%v]", data.MessageID, items[0].GetMessageID())
	}
}
