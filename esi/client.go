package esi

import (
	"encoding/json"
	"fmt"
	"github.com/dariusbakunas/eve-processors/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	BaseUrl string
	Token   string
	Timeout time.Duration
}

func NewEsiClient(baseUrl string, token string, timeout time.Duration) *Client {
	return &Client{
		BaseUrl: baseUrl,
		Token:   token,
		Timeout: timeout,
	}
}

func (c *Client) get(url string) ([]byte, *http.Header, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	errorLimitRemainStr := resp.Header.Get("X-ESI-Error-Limit-Remain")
	errorLimitRemain, err := strconv.Atoi(errorLimitRemainStr)

	if err != nil {
		return nil, &resp.Header, fmt.Errorf("strconv.Atoi: %v", err)
	}

	log.Printf("X-ESI-Error-Limit-Remain: %d", errorLimitRemain)

	errorLimitResetStr := resp.Header.Get("X-ESI-Error-Limit-Reset")
	errorLimitReset, err := strconv.Atoi(errorLimitResetStr)

	if err != nil {
		return nil, &resp.Header, fmt.Errorf("strconv.Atoi: %v", err)
	}

	log.Printf("X-ESI-Error-Limit-Reset: %d", errorLimitReset)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &resp.Header, err
	}

	if 200 != resp.StatusCode {
		return nil, &resp.Header, fmt.Errorf("%s", body)
	}

	return body, &resp.Header, nil
}

type JournalEntriesResponse struct {
	Pages   int
	Entries []models.JournalEntry
}

func (c *Client) GetWalletTransactions(characterId int64) ([]models.WalletTransaction, error) {
	url := fmt.Sprintf("%s/characters/%d/wallet/transactions/", c.BaseUrl, characterId)
	bytes, _, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var data []models.WalletTransaction
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetJournalEntries(characterId int64, page int) (*JournalEntriesResponse, error) {
	url := fmt.Sprintf("%s/characters/%d/wallet/journal/?page=%d", c.BaseUrl, characterId, page)
	bytes, headers, err := c.get(url)
	if err != nil {
		return nil, err
	}

	pagesStr := headers.Get("X-Pages")

	pages, err := strconv.Atoi(pagesStr)

	if err != nil {
		return nil, fmt.Errorf("strconv.Atoi: %v", err)
	}

	var data []models.JournalEntry
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &JournalEntriesResponse{
		Pages:   pages,
		Entries: data,
	}, nil
}
