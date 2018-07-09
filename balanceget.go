package twizo

import (
	"encoding/json"
	"net/http"
)

type jsonBalanceGetResponse struct {
	Credit            float32 `json:"credit"`
	CurrencyCode      string  `json:"currencyCode"`
	Wallet            string  `json:"wallet"`
	AlarmLimit        *string `json:"alarmLimit,omitempty"`
	FreeVerifications int     `json:"freeVerifications,omitempty"`
}

// BalanceGetResponse struct that the server returns for a verification request
type BalanceGetResponse struct {
	credit            float32
	currencyCode      string
	wallet            string
	alarmLimit        *string
	freeVerifications int
}

// UnmarshalJSON unmarshals the json returned by the server into the BalanceGetResponse struct
func (response *BalanceGetResponse) UnmarshalJSON(j []byte) error {
	var jsonResponse = &jsonBalanceGetResponse{}

	err := json.Unmarshal(j, &jsonResponse)
	if err != nil {
		return err
	}

	return response.copyFrom(jsonResponse)
}

func (response *BalanceGetResponse) copyFrom(j *jsonBalanceGetResponse) error {
	var err error // default err is nil

	response.credit = j.Credit
	response.currencyCode = j.CurrencyCode
	response.wallet = j.Wallet
	response.alarmLimit = j.AlarmLimit
	response.freeVerifications = j.FreeVerifications

	return err
}

// GetCredit get the current credit
func (response BalanceGetResponse) GetCredit() float32 {
	return response.credit
}

// GetCurrencyCode get the current credit currency code
func (response BalanceGetResponse) GetCurrencyCode() string {
	return response.currencyCode
}

// GetWallet gets the wallet name
func (response BalanceGetResponse) GetWallet() string {
	return response.wallet
}

// GetAlarmLimit get the current alarm limit
func (response BalanceGetResponse) GetAlarmLimit() *string {
	return response.alarmLimit
}

// GetFreeVerifications get the amount of free verifications left
func (response BalanceGetResponse) GetFreeVerifications() int {
	return response.freeVerifications
}

// BalanceGetRequest empty struct placeholder (future use)
type BalanceGetRequest struct{}

// Submit wil submit the balance request
func (request *BalanceGetRequest) Submit() (*BalanceGetResponse, error) {
	response := &BalanceGetResponse{}

	apiURL, err := GetURLFor("wallet/getbalance")
	if err != nil {
		return nil, err
	}

	err = GetClient(RegionCurrent, APIKey).Call(
		http.MethodGet,
		apiURL,
		nil,
		http.StatusOK,
		response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// NewBalanceGetRequest creates a new BalanceRequest
func NewBalanceGetRequest() *BalanceGetRequest {
	request := &BalanceGetRequest{}
	return request
}

// BalanceGet retrieves the credit balance of the api key
func BalanceGet() (*BalanceGetResponse, error) {
	request := NewBalanceGetRequest()
	return request.Submit()
}
