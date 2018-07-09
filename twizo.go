package twizo

// issues ?
//
// Api response difference ? add error detail codes ? we need to string compare now
// - {"title":"Not Found","status":404,"detail":"Entity not found."}
// - {"title":"Not Found","status":404,"detail":"Page not found."}

// Consistency on _embedded
// - SmsSubmit + NumberLookupSubmit returns _embedded  -> items,
// - PollSubmit returns embedded in _embedded -> messages

// Consistency on count / total_items
// - SmsSubmit + NumberLookupSubmit returns total_items
// - PollSubit returns count

// Documentation:
// - Add the types explicitly .. ie is it always a string / int / []string etc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// APIRegion the region to call
type APIRegion string

// ResultType how to we want the results
type ResultType int

// All api regions
const (
	APIRegionAsia    APIRegion = "asia"
	APIRegionEU      APIRegion = "eu"
	APIRegionDefault APIRegion = "default"
)

// All client variables
const (
	ClientAPIVersion string = "v1"
	ClientLibVersion string = "0.1.0"
	ClientLibName    string = "Twizo-go-lib"
	ClientAuthUser   string = "twizo"
)

const (
	defaultHTTPTimeout = 80 * time.Second
)

// All result types
const (
	// ResultTypeNone (default)
	ResultTypeNone ResultType = 0

	// ResultTypeCallback send result to callbackURL (via post)
	ResultTypeCallback ResultType = 1

	// ResultTypePolling we want to poll for the results
	ResultTypePolling ResultType = 2

	// ResultTypeCallbackPolling (use both SmsResultTypeCallback
	// and SmsResultTypePolling)
	ResultTypeCallbackPolling ResultType = 3
)

// DebugLogger contains the debug logger
var DebugLogger = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

var regionUrls = map[APIRegion]string{
	APIRegionAsia:    "api-asia-01.twizo.com",
	APIRegionEU:      "api-eu-01.twizo.com",
	APIRegionDefault: "api-eu-01.twizo.com",
}

var httpClient = &http.Client{Timeout: HTTPClientTimeout}
var httpClientUserAgent = fmt.Sprintf(
	"%s/%s Go/%s/%s/%s",
	ClientLibName,
	ClientLibVersion,
	runtime.Version(),
	runtime.GOARCH,
	runtime.GOOS,
)

// APIKey the api key : possible to override this setting
var APIKey string

// RegionCurrent the current region : possible to override this setting
var RegionCurrent = APIRegionDefault

// HTTPClientTimeout the current timeout on http calls
var HTTPClientTimeout = defaultHTTPTimeout

// Recipient the recipient
type Recipient string

// Request interface
type Request interface {
}

// Response interface
type Response interface {
	UnmarshalJSON(data []byte) error
}

// Client interface
type Client interface {
	Call(method string, path string, request Request, v interface{}) error
}

// HTTPClient is the actual http client
type HTTPClient struct {
	Region     APIRegion
	Key        string
	HTTPClient *http.Client
}

// NewRequest creates a new request, this allows it to be tested [todo: refactor]
func (c HTTPClient) NewRequest(method string, url *url.URL, body io.Reader) (*http.Request, error) {

	if url.Host == "" {
		url.Host = GetHostForRegion(c.Region)
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(ClientAuthUser, c.Key)

	req.Header.Add("User-Agent", httpClientUserAgent)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// Call performs the actual call on the client
func (c *HTTPClient) Call(method string, url *url.URL, request Request, expectCode int, v interface{}) error {
	// convert request to body
	requestBody := bytes.NewBuffer(nil)
	if request != nil {
		body, err := json.Marshal(request)
		if err != nil {
			return err
		}

		requestBody.Write([]byte(body))
	}

	// create new request
	req, err := c.NewRequest(method, url, requestBody)
	if err != nil {
		return err
	}

	if requestBody.Len() > 0 {
		DebugLogger.Printf("Request %v [%v] body %s", req.Method, req.URL.String(), requestBody)
	} else {
		DebugLogger.Printf("Request %v [%v]", req.Method, req.URL.String())
	}

	// actually do the request and parse errors if any
	if err := c.do(req, expectCode, v); err != nil {
		return err
	}

	return nil
}

// Do is used by Call to execute an API request and parse the response. It uses
// the backend's HTTP client to execute the request and unmarshals the response
// into v. It also handles unmarshaling errors returned by the API.
func (c *HTTPClient) do(request *http.Request, expectCode int, response interface{}) error {
	start := time.Now()

	res, err := c.HTTPClient.Do(request)

	if err != nil {
		return err
	}
	defer res.Body.Close() // nolint: errcheck

	// might want to use json.Decoder instead of ioutl.ReadAll -> sending to Unmarshal
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(resBody) > 0 {
		DebugLogger.Printf("Response in [%v] with [%d] body %s", time.Since(start), res.StatusCode, resBody)
	} else {
		DebugLogger.Printf("Response in [%v] with [%d]", time.Since(start), res.StatusCode)
	}

	// check if there was a problem
	if res.Header.Get("Content-Type") == "application/problem+json" {
		if res.StatusCode == http.StatusUnprocessableEntity {
			apiError := &APIValidationError{}
			if err := json.Unmarshal(resBody, apiError); err != nil {
				return err
			}
			return apiError
		}

		apiError := &APIError{}
		if err := json.Unmarshal(resBody, apiError); err != nil {
			return err
		}
		return apiError
	}

	// check if the response code was something we expected
	if res.StatusCode != expectCode {
		clientError := &ClientError{
			Message: fmt.Sprintf("Unexpected response [%s]", resBody),
			Code:    res.StatusCode,
		}
		return clientError
	}

	// todo check if we are actually a http NoContent here ?
	if response == nil {
		// not expecting result we are done
		return nil
	}

	// https://ahmetalpbalkan.com/blog/golang-json-decoder-pitfalls/
	// todo: check the content-type to make sure it's application/json
	if err := json.Unmarshal(resBody, response); err != nil {
		return err
	}

	return nil
}

// GetClient gets the a client initialized with region and key
func GetClient(region APIRegion, key string) *HTTPClient {
	return &HTTPClient{region, key, GetHTTPClient()}
}

//
// Support Structs
//

// HATEOASHref link structure
type HATEOASHref struct {
	Href url.URL `json:"href"`
}

// UnmarshalJSON the HATEOASHref link
func (l *HATEOASHref) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "href" {
			u, err := url.Parse(v)
			if err != nil {
				return err
			}
			l.Href = *u
			break
		}
	}

	return nil
}

// HATEOASLinks struct
type HATEOASLinks struct {
	Self HATEOASHref `json:"self"`
}

func (h HATEOASLinks) getDeepClone() HATEOASLinks {
	linkURL, _ := url.Parse(h.Self.Href.String())

	return HATEOASLinks{
		Self: HATEOASHref{
			Href: *linkURL,
		},
	}
}

func createSelfLinks(selfLink *url.URL) HATEOASLinks {
	links := HATEOASLinks{Self: HATEOASHref{Href: *selfLink}}
	return links
}

//
// Support functions
//

// SetHTTPClient overrides the default HTTP client.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func SetHTTPClient(client *http.Client) {
	httpClient = client
}

// GetHTTPClient returns the current http client
func GetHTTPClient() *http.Client {
	return httpClient
}

// GetURLFor allows mocking
var GetURLFor = GetURLForOriginal

// GetURLForOriginal gets the url for a request [todo: refactor]
func GetURLForOriginal(path string) (*url.URL, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if !strings.HasPrefix(path, "/"+ClientAPIVersion) {
		path = fmt.Sprintf("/%s%s", ClientAPIVersion, path)
	}

	apiURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// always force https
	apiURL.Scheme = "https"

	return apiURL, nil
}

// AddHostForRegion adds a region with a host
func AddHostForRegion(region APIRegion, host string) {
	regionUrls[region] = host
}

// GetHostForRegion gets a region for a host
func GetHostForRegion(region APIRegion) string {
	if host, ok := regionUrls[region]; ok {
		return host
	}

	return regionUrls["default"]
}

// GetRegions get all regions
func GetRegions() map[APIRegion]string {
	return regionUrls
}

func isDcsBinary(i int) bool {
	if i&200 == 0 || i&248 == 240 {
		return (i & 4) > 0
	}
	return false
}

func convertRecipients(recipients interface{}) ([]Recipient, error) {
	var r []Recipient
	switch tRecipients := recipients.(type) {
	case []Recipient:
		r = tRecipients
	case Recipient:
		r = []Recipient{tRecipients}
	case string:
		r = []Recipient{Recipient(tRecipients)}
	case []string:
		for _, element := range tRecipients {
			r = append(r, Recipient(element))
		}
	default:
		return nil, fmt.Errorf("expecting string []Recipient or Recipient, got [%s]", reflect.TypeOf(recipients))
	}
	return r, nil
}
