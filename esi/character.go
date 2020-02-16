package esi

import (
	"encoding/json"
	"fmt"
	"github.com/dariusbakunas/eve-processors/models"
	"strconv"
)

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

func (c *Client) GetCharacterMarketOrders(characterID int64) ([]models.CharacterMarketOrder, error) {
	url := fmt.Sprintf("%s/characters/%d/orders", c.BaseUrl, characterID)
	bytes, _, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("c.get: %v", err)
	}
	var data []models.CharacterMarketOrder
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

	var data []models.CharacterMarketOrder
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &models.MarketOrderHistoryResponse{
		Orders: data,
		Pages:  pages,
	}, nil
}

func (c *Client) GetBlueprints(characterID int64, page int) (*models.BlueprintsResponse, error) {
	url := fmt.Sprintf("%s/characters/%d/blueprints/?page=%d", c.BaseUrl, characterID, page)
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

	var data []models.Blueprint
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &models.BlueprintsResponse{
		Blueprints: data,
		Pages:  pages,
	}, nil
}

func (c *Client) GetIndustryJobs(characterID int64) ([]models.IndustryJob, error) {
	url := fmt.Sprintf("%s/characters/%d/industry/jobs/?include_completed=true", c.BaseUrl, characterID)

	bytes, _, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("c.get: %v", err)
	}

	var data []models.IndustryJob
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return data, nil
}