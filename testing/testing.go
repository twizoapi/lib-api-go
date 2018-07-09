package testing

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	twizo "github.com/twizoapi/lib-api-go"
	"gopkg.in/jarcoal/httpmock.v1"
)

// Constants in use for testing
const (
	TestAPIKey = "test-api-key"
	TestRegion = "eu"
)

// ParseTemplateStringToBytes parses a template, data and returns the bytes
func ParseTemplateStringToBytes(tmplStr string, data interface{}) ([]byte, error) {
	templateNew, err := template.New("template").Parse(tmplStr)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	bWriter := bufio.NewWriter(&b)
	err = templateNew.Execute(bWriter, data)
	if err != nil {
		return nil, err
	}
	bWriter.Flush()

	return b.Bytes(), nil
}

// ValidateRequestFunc - testing -
type ValidateRequestFunc func(req *http.Request) error

// HTTPMockSend mocks a http get
func HTTPMockSend(Method string, URL string, expect int, response interface{}, vFunc ValidateRequestFunc) error {
	resp := new(http.Response)

	if vFunc == nil {
		vFunc = func(req *http.Request) error { return nil }
	}

	if response == nil && expect != http.StatusNoContent {
		return fmt.Errorf("httpmocksend: no nil response expected")
	}

	if expect == http.StatusNoContent && response != nil {
		return fmt.Errorf("httpmocksend: unknown response item: %#v for NoContent", response)
	}

	switch v := response.(type) {
	case *twizo.APIError:
		const tpl = `{
				"type":"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
				"title":"{{.Title}}",
				"status":{{.Status}},
				"detail":"{{.Detail}}"
			}`
		data := struct {
			Title  string
			Status int
			Detail string
		}{
			Title:  v.Title(),
			Status: v.Status(),
			Detail: v.Detail(),
		}
		b, err := ParseTemplateStringToBytes(tpl, data)
		if err != nil {
			return err
		}
		resp = httpmock.NewBytesResponse(expect, b)
		resp.Header.Set("Content-Type", "application/problem+json")
	case string:
		resp = httpmock.NewStringResponse(expect, v)
		resp.Header.Set("Content-Type", "application/json")
	case []byte:
		var TestResponse interface{}
		err := json.Unmarshal(v, &TestResponse)
		if err != nil {
			return err
		}
		resp, err = httpmock.NewJsonResponse(expect, TestResponse)
		if err != nil {
			resp = httpmock.NewStringResponse(http.StatusInternalServerError, err.Error())
		}
	case nil:
		resp = httpmock.NewStringResponse(expect, "")
	default:
		return fmt.Errorf("httpmocksend: unknown response item: %#v", response)
	}

	httpmock.RegisterResponder(Method,
		URL,
		func(req *http.Request) (*http.Response, error) {
			err := vFunc(req)
			if err != nil {
				return nil, err
			}
			return resp, nil
		},
	)
	return nil
}
