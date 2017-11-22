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

// NumberLookupPollResults results of the number lookup poll
type NumberLookupPollResults struct {
	pollResults
	Embedded *numberLookupPollEmbedded `json:"_embedded"`
}

func (p NumberLookupPollResults) getURL() (*url.URL, error) {
	return GetURLFor("numberlookup/poll")
}

// GetItems returns the embedded messages of a poll action
func (p NumberLookupPollResults) GetItems() []NumberLookupResponse {
	return *p.Embedded.Messages
}

// NewNumberLookupPoll creates a new NumberLookupPoll
func NewNumberLookupPoll() *NumberLookupPollResults {
	params := &NumberLookupPollResults{}
	return params
}

// NumberLookupPollStatus [todo: refactor]
func NumberLookupPollStatus() (*NumberLookupPollResults, error) {
	request := NewNumberLookupPoll()
	err := request.Status()

	return request, err
}

// Status requests the status of a NumberLookupPollResults
func (p *NumberLookupPollResults) Status() error {
	apiURL, err := p.getURL()
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
