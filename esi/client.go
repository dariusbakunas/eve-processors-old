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

import retry "github.com/hashicorp/go-retryablehttp"

type Client struct {
	BaseUrl string
	Token   string
	HttpClient *retry.Client
}

func NewEsiClient(baseUrl string, token string, timeout time.Duration) *Client {
	client := retry.NewClient()

	return &Client{
		BaseUrl: baseUrl,
		Token:   token,
		HttpClient: client,
	}
}

func (c *Client) get(url string) ([]byte, *http.Header, error) {
	req, err := retry.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	client := c.HttpClient
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

func (c *Client) GetJournalEntries(characterId int64, page int) (*models.JournalEntriesResponse, error) {
	url := fmt.Sprintf("%s/characters/%d/wallet/journal/?page=%d", c.BaseUrl, characterId, page)
	bytes, headers, err := c.get(url)
	if err != nil {
		return nil, err
	}

	pagesStr := headers.Get("X-Pages")

	pages := 1

	if pagesStr != "" {
		pages, err = strconv.Atoi(pagesStr)
		if err != nil {
			return nil, fmt.Errorf("strconv.Atoi: %v", err)
		}
	}

	var data []models.JournalEntry
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return &models.JournalEntriesResponse{
		Pages:   pages,
		Entries: data,
	}, nil
}

func (c *Client) GetSkills(characterID int64) (*models.SkillsResponse, error) {
	url := fmt.Sprintf("%s/characters/%d/skills", c.BaseUrl, characterID)
	bytes, _, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("c.get: %v", err)
	}

	var data models.SkillsResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return &data, nil
}

func (c *Client) GetSkillQueue(characterID int64) ([]models.SkillQueueItem, error) {
	url := fmt.Sprintf("%s/characters/%d/skillqueue", c.BaseUrl, characterID)
	bytes, _, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("c.get: %v", err)
	}
	var data []models.SkillQueueItem
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetMarketOrders(characterID int64) ([]models.MarketOrder, error) {
	url := fmt.Sprintf("%s/characters/%d/orders", c.BaseUrl, characterID)
	bytes, _, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("c.get: %v", err)
	}
	var data []models.MarketOrder
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetMarketOrderHistory(characterID int64, page int) (*models.MarketOrderHistoryResponse, error) {
	url := fmt.Sprintf("%s/characters/%d/orders/history/?page=%d", c.BaseUrl, characterID, page)
	bytes, headers, err := c.get(url)
	if err != nil {
		return nil, err
	}

	pagesStr := headers.Get("X-Pages")

	pages := 1

	if pagesStr != "" {
		pages, err = strconv.Atoi(pagesStr)
		if err != nil {
			return nil, fmt.Errorf("strconv.Atoi: %v", err)
		}
	}

	var data []models.MarketOrder
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &models.MarketOrderHistoryResponse{
		Orders: data,
		Pages:  pages,
	}, nil
}
