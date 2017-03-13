package twizo

import (
	"net/url"
)

//
// Polling support
//
type numberLookupPollEmbedded struct {
	Messages *[]NumberLookupResponse `json:"messages"` // strange this was called items in smsResponseEmbedded
}

type NumberLookupPollResults struct {
	pollResults
	Embedded *numberLookupPollEmbedded `json:"_embedded"`
}

func (p *NumberLookupPollResults) GetURL() (*url.URL, error) {
	return GetURLFor("numberlookup/poll")
}

func (p *NumberLookupPollResults) GetItems() []NumberLookupResponse {
	return *p.Embedded.Messages
}

func NewNumberLookupPoll() *NumberLookupPollResults {
	params := &NumberLookupPollResults{}
	return params
}

func NumberLookupPollStatus() (*NumberLookupPollResults, error) {
	request := NewNumberLookupPoll()
	err := request.Status()

	return request, err
}

func (p *NumberLookupPollResults) Status() error {
	apiURL, err := p.GetURL()
	if err != nil {
		return err
	}
	resp := NewNumberLookupPoll()
	err = p.status(apiURL, resp)
	if err == nil {
		*p = *resp
	}
	return err
}
