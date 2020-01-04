package esi

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"time"
)

type EsiClient struct {
	BaseUrl string
	Token string
	Timeout time.Duration
}

func NewEsiClient(baseUrl string, token string, timeout time.Duration) *EsiClient {
	return &EsiClient{
		BaseUrl: baseUrl,
		Token:   token,
		Timeout: timeout,
	}
}

func (c *EsiClient) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

type WalletTransaction struct {
	ClientId      uint            `json:"client_id"`
	Quantity      uint            `json:"quantity"`
	UnitPrice     decimal.Decimal `json:"unit_price"`
	Date          time.Time       `json:"date"`
	IsBuy         bool            `json:"is_buy"`
	IsPersonal    bool            `json:"is_personal"`
	JournalRefId  uint64          `json:"journal_ref_id"`
	LocationId    uint64          `json:"location_id"`
	TransactionId uint64          `json:"transaction_id"`
	TypeId        uint            `json:"type_id"`
}

func (c *EsiClient) GetWalletTransactions(characterId int64) ([]WalletTransaction, error) {
	url := fmt.Sprintf("%s/characters/%d/wallet/transactions/", c.BaseUrl, characterId)
	bytes, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var data []WalletTransaction
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}