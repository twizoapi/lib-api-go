package twizo_test

import (
	"encoding/json"
	"reflect"
	"testing"

	twizo "github.com/twizoapi/lib-api-go"
	. "github.com/twizoapi/lib-api-go/testing"
)

func init() {
	twizo.APIKey = TestAPIKey
	twizo.RegionCurrent = TestRegion
}

func TestWidgetSessionInvalidJsonResponse(t *testing.T) {
	jsonResponse := &twizo.WidgetSessionResponse{}
	err := jsonResponse.UnmarshalJSON([]byte("Invalid json"))
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf(
			"Invalid error expecting [json.SyntaxError] got [%#v]",
			err,
		)
	}
}

func TestWidgetSessionRequest(t *testing.T) {
	data := struct {
		Issuer string
	}{
		Issuer: "issuer",
	}

	sessionRequest := twizo.NewWidgetSessionRequest()
	sessionRequest.SetIssuer(data.Issuer)
	if data.Issuer != sessionRequest.GetIssuer() {
		t.Fatalf(
			"Invalid Issuer expecting [%s] got [%s]",
			data.Issuer,
			sessionRequest.GetIssuer(),
		)
	}

	// please not that Order is arrays must be maintaned as the 1st has perference over the 2nd
	sessionRequest.SetAllowedTypes([]string{"line", "sms"})
	if !reflect.DeepEqual(sessionRequest.GetAllowedTypes(), twizo.VerificationTypes{twizo.VerificationTypeLine, twizo.VerificationTypeSms}) {
		t.Fatalf(
			"Slices should be the same expected [%#v] got [%#v]",
			twizo.VerificationTypes{twizo.VerificationTypeSms, twizo.VerificationTypeLine},
			sessionRequest.GetAllowedTypes(),
		)
	}
}
