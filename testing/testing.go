package testing

import (
	"text/template"
	"bytes"
	"bufio"
	"encoding/json"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
)

const (
	TestApiKey = "test-api-key"
	TestRegion = "eu"
)

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


func HttpMockSendJsonPostTo(URL string, expect int, b []byte) (error) {
	var TestResponse interface{}
	err := json.Unmarshal(b, &TestResponse)
	if err != nil {
		return err
	}
	httpmock.RegisterResponder(http.MethodPost,
		URL,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(expect, TestResponse)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	return nil
}
