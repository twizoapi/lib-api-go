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

func TestSmsNew(t *testing.T) {
	smsRequest := twizo.NewSmsRequest([]twizo.Recipient{twizo.Recipient("0000000000")}, "Message", "Sender")

	_, err := json.Marshal(smsRequest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSmsSubmit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// our basic response we will need to add more here
	const tpl = `{
		"_links": {
			"self":{"href":"https:\/\/{{ .Host }}\/v1\/sms\/submitsimple"}
		},
		"_embedded":{
			"items":[
				{
					"applicationTag":"Test application",
					"body":"Message",
					"callbackUrl":null,
					"createdDateTime":"2017-03-09T14:47:49+00:00",
					"dcs":0,
					"messageId":"{{ .MessageID }}",
					"networkCode":null,
					"pid":null,
					"reasonCode":null,
					"recipient":"{{ .Number }}",
					"resultTimestamp":null,
					"resultType":0,
					"salesPrice":null,
					"salesPriceCurrencyCode":null,
					"scheduledDelivery":null,
					"sender":"TwizoTest",
					"senderNpi":0,
					"senderTon":5,
					"status":"no status",
					"statusCode":0,
					"tag":null,
					"udh":null,
					"validity":259200,
					"validUntilDateTime":"2017-03-12T14:47:49+00:00",
					"_links":{
						"self":{"href":"https:\/\/{{ .Host }}\/v1\/sms\/submitsimple\/{{ .MessageID }}"}
					}
				}
			]
		},
		"total_items":1
	}`
	data := struct {
		MessageID string
		Number    string
		Host      string
	} {
		MessageID : "test-1.10314.sms58c16b15c261a5.18930279",
		Number    : "6100000000",
		Host      : twizo.GetHostForRegion(twizo.RegionCurrent),
	}

	b, err := ParseTemplateStringToBytes(tpl, data)
	if err != nil {
		t.Fatal(err)
	}

	HttpMockSendJsonPostTo(
		fmt.Sprintf("https://%s/%s/sms/submitsimple", data.Host, twizo.ClientAPIVersion),
		http.StatusCreated,
		b,
	)

	response, err := twizo.SmsSubmit([]twizo.Recipient{twizo.Recipient(data.Number)}, "Message", "TwizoTest")
	if err != nil {
		t.Fatal(err)
	}
	items := response.GetItems()

	// only check for now check messageId
	if items[0].GetMessageID() != data.MessageID {
		t.Fatalf("Invalid message id expected [%s] got [%v]", data.MessageID, items[0].GetMessageID())
	}
}

func TestSmsBinary(t *testing.T) {
	smsRequest := twizo.SmsRequest{}
	smsRequest.SetDCS(0xF5)

	if !smsRequest.IsBinary() {
		t.Fatal("SmsRequest should be binary it's not")
	}

	smsRequest = twizo.SmsRequest{}
	if smsRequest.IsBinary() {
		t.Fatal("SmsRequest should not be binary it is")
	}
}
