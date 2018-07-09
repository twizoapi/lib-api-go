package twizo

import (
	"net/http"
	"net/url"
)

type pollResultInterface interface {
	GetURL() (*url.URL, error)
}

type pollResults struct {
	pollResultInterface
	BatchID string        `json:"batchId"`
	Count   int           `json:"count"` // strange was called total_items in smsResponse
	Links   *HATEOASLinks `json:"_links"`
}

func (p *pollResults) status(url *url.URL, resp interface{}) error {
	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		url,
		nil,
		http.StatusOK,
		resp,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *pollResults) Delete() error {
	if p.BatchID == "" {
		// we can not delete without a batchId,
		return nil
	}

	// we can delete the same thing over and over and still get the
	// statusNoContent response, so this is safe to do.
	err := GetClient(RegionCurrent, APIKey).Call(
		http.MethodDelete,
		&p.Links.Self.Href,
		nil,
		http.StatusNoContent,
		nil,
	)

	return err
}
