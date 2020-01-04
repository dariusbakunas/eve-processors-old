package esi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetTokens(clientId string, clientSecret string, refreshToken string) (*TokensResponse, error) {
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientSecret)))
	url := "https://login.eveonline.com/oauth/token"
	jsonPayload := []byte(fmt.Sprintf(`{"grant_type":"refresh_token", "refresh_token": "%s"}`, refreshToken))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 3,
	}

	log.Printf("Sending request to %s", url)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	if 200 != resp.StatusCode {
		for k, v := range resp.Header {
			log.Printf("%s : %s", k, v)
		}
		return nil, fmt.Errorf("resp.StatusCode: %s", body)
	}

	var data TokensResponse
	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &data, nil
}

type TokensResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int64 `json:"expires_in"`
}

func getCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetAccessToken(db *DB, character Character, eveClientId string, eveClientSecret string) (string, error) {
	timestamp := getCurrentTimestamp()

	if timestamp > int64(character.expires-1000*60) {
		log.Printf("Updating tokens for Character ID: %d", character.id)

		tokens, err := GetTokens(eveClientId, eveClientSecret, character.refreshToken)

		if err != nil {
			return "", fmt.Errorf("getTokens: %v", err)
		}

		err = db.UpdateCharacterTokens(tokens.AccessToken, tokens.RefreshToken, tokens.ExpiresIn, character.id)

		if err != nil {
			return "", fmt.Errorf("sq.Update: %v", err)
		}

		return tokens.AccessToken, nil
	} else {
		return character.accessToken, nil
	}
}
