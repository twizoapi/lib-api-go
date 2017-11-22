package twizo

import (
	"net/url"
)

//
// Polling support
//

type smsPollEmbedded struct {
	Messages *[]SmsResponse `json:"messages"` // strange this was called items in smsResponseEmbedded
}

// SmsPollResults struct
type SmsPollResults struct {
	pollResults                  // an anonymous field of type pollResult
	Embedded    *smsPollEmbedded `json:"_embedded"`
}

func (p SmsPollResults) getURL() (*url.URL, error) {
	return GetURLFor("sms/poll")
}

// GetItems get the embedded items of a poll result
func (p SmsPollResults) GetItems() []SmsResponse {
	return *p.Embedded.Messages
}

// NewSmsPoll creates a new smspollresult
func NewSmsPoll() *SmsPollResults {
	params := &SmsPollResults{}
	return params
}

// SmsPollStatus [todo: refactor]
func SmsPollStatus() (*SmsPollResults, error) {
	request := NewSmsPoll()
	err := request.Status()

	return request, err
}

// Status get the status of a sms poll result
func (p *SmsPollResults) Status() error {
	apiURL, err := p.getURL()
	if err != nil {
		return err
	}
	resp := NewSmsPoll()
	err = p.status(apiURL, resp)
	if err == nil {
		*p = *resp
	}
	return err
}
