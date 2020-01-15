package esi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dariusbakunas/eve-processors/models"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

import retry "github.com/hashicorp/go-retryablehttp"

type Client struct {
	BaseUrl string
	Token   null.String
	HttpClient *retry.Client
}

func NewAuthenticatedEsiClient(baseUrl string, token string, timeout time.Duration) *Client {
	client := retry.NewClient()

	return &Client{
		BaseUrl: baseUrl,
		Token:   null.NewString(token, true),
		HttpClient: client,
	}
}

func NewEsiClient(baseUrl string) *Client {
	client := retry.NewClient()

	return &Client{
		BaseUrl: baseUrl,
		Token:   null.String{},
		HttpClient: client,
	}
}

type pagedGetResult struct {
	index int
	data  []byte
	err   error
}

func (c *Client) pagedGet(urls []string, concurrencyLimit int) []pagedGetResult {
	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrencyLimit)
	stop := make(chan struct{})

	// this channel will not block and collect the http request results
	resultsChan := make(chan *pagedGetResult)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i, url := range urls {
		go func(i int, url string) {
			// do not continue if stop channel is closed
			select {
			default:
				semaphoreChan <- struct{}{}
				req, err := retry.NewRequest("GET", url, nil)

				if err != nil {
					result := &pagedGetResult{i, nil, err}
					resultsChan <- result
					return
				}

				if c.Token.Valid {
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.String))
				}

				req.Header.Set("Content-Type", "application/json")

				client := c.HttpClient

				resp, err := client.Do(req)

				if err != nil {
					result := &pagedGetResult{i, nil, err}
					resultsChan <- result
					return
				}

				errorLimitRemainStr := resp.Header.Get("X-ESI-Error-Limit-Remain")
				errorLimitRemain, err := strconv.Atoi(errorLimitRemainStr)

				if err != nil {
					result := &pagedGetResult{i, nil, err}
					resultsChan <- result
					return
				}

				if errorLimitRemain == 0 {
					// close stop chanel to prevent remaining routines from executing
					log.Printf("Esi error limit reached, cancelling all requests")
					result := &pagedGetResult{i, nil, errors.New("esi error limit reached, cancelling all requests")}
					resultsChan <- result
					close(stop)
					return
				}

				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					result := &pagedGetResult{i, nil, err}
					resultsChan <- result
					return
				}

				if 200 != resp.StatusCode {
					result := &pagedGetResult{i, nil, fmt.Errorf("%s", body)}
					resultsChan <- result
					return
				}

				result := &pagedGetResult{i, body, nil}
				resultsChan <- result

				// once we're done it's we read from the semaphoreChan which
				// has the effect of removing one from the limit and allowing
				// another goroutine to start
				<-semaphoreChan
			case <-stop: // triggered when the stop channel is closed
				result := &pagedGetResult{i, nil, errors.New("execution cancelled due to esi error limit")}
				resultsChan <- result
				<-semaphoreChan
				return
			}
		}(i, url)
	}

	var results []pagedGetResult

	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(urls) {
			break
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	return results
}

func (c *Client) get(url string) ([]byte, *http.Header, error) {
	req, err := retry.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	if c.Token.Valid {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.String))
	}

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

func (c *Client) GetMarketOrders(regionID int) (map[int][]models.MarketOrder, error) {
	url := fmt.Sprintf("%s/markets/%d/orders/?page=%d", c.BaseUrl, regionID, 1)
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

	var marketOrders []models.MarketOrder

	err = json.Unmarshal(bytes, &marketOrders)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	if pages > 1 {
		urls := make([]string, pages - 1)

		for i := range urls {
			urls[i] = fmt.Sprintf("%s/markets/%d/orders/?page=%d", c.BaseUrl, regionID, i + 2)
		}

		pagedResults := c.pagedGet(urls, 10)

		for _, result := range pagedResults {
			var orders []models.MarketOrder

			if result.data != nil {
				err = json.Unmarshal(result.data, &orders)
				if err != nil {
					return nil, fmt.Errorf("json.Unmarshal: %v", err)
				}

				marketOrders = append(marketOrders, orders...)
			} else {
				// not returning error, using partial results
				// TODO: maybe make this conditional?
				log.Printf("Failed to get market orders, page: %d, error: %v", result.index + 2, result.err)
			}
		}
	}

	orderMap := make(map[int][]models.MarketOrder)

	for _, order := range marketOrders {
		if _, found := orderMap[order.TypeID]; found {
			orderMap[order.TypeID] = append(orderMap[order.TypeID], order)
		} else {
			orderMap[order.TypeID] = []models.MarketOrder{order}
		}
	}

	return orderMap, nil
}
