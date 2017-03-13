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

type SmsPollResults struct {
	pollResults             // an anonymous field of type pollResult
	Embedded   *smsPollEmbedded `json:"_embedded"`
}

func (p *SmsPollResults) GetURL() (*url.URL, error) {
	return GetURLFor("sms/poll")
}

func (p *SmsPollResults) GetItems() []SmsResponse {
	return *p.Embedded.Messages
}

func NewSmsPoll() *SmsPollResults {
	params := &SmsPollResults{}
	return params
}

func SmsPollStatus() (*SmsPollResults, error) {
	request := NewSmsPoll()
	err := request.Status()

	return request, err
}

func (p *SmsPollResults) Status() error {
	apiURL, err := p.GetURL()
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
