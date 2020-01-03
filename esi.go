package esi

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

import _ "github.com/go-sql-driver/mysql"
import sq "github.com/Masterminds/squirrel"

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func initializeDb() (*sql.DB, error) {
	dbConnection := os.Getenv("DB_CONNECTION")

	if dbConnection == "" {
		return nil, fmt.Errorf("DB_CONNECTION must be set")
	}

	dbDatabase := os.Getenv("DB_DATABASE")

	if dbDatabase == "" {
		return nil, fmt.Errorf("DB_DATABASE must be set")
	}

	dbUsername := os.Getenv("DB_USERNAME")

	if dbUsername == "" {
		return nil, fmt.Errorf("DB_USERNAME must be set")
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD must be set")
	}

	var err error
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", dbUsername, dbPassword, dbConnection, dbDatabase))

	if err != nil {
		return nil, fmt.Errorf("mysql: %v", err)
	}

	return db, nil
}

type Character struct {
	id           int
	accessToken  string
	refreshToken string
	expires      int
	scopes       string
}

func getCharacters(db *sql.DB) ([]Character, error) {
	rows, err := sq.Select("id, accessToken, refreshToken, expiresAt, scopes").From("characters").RunWith(db).Query()

	if err != nil {
		return nil, fmt.Errorf("mysql: %v", err)
	}

	defer rows.Close()

	var characters []Character

	for rows.Next() {
		var character Character

		err := rows.Scan(&character.id, &character.accessToken, &character.refreshToken, &character.expires, &character.scopes)

		if err != nil {
			return nil, fmt.Errorf("mysql: %v", err)
		}

		characters = append(characters, character)
	}

	return characters, nil
}

func processWalletTransactions(client *EsiClient, characterId int, pubSub bool) error {
	if pubSub {
		log.Printf("Sending message")
		return nil
	}

	transactions, err := client.GetWalletTransactions(characterId)

	if err != nil {
		return err
	}

	fmt.Printf("%+v", transactions)

	return nil
}

func getCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type TokensResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int64 `json:"expires_in"`
}

func getTokens(clientId string, clientSecret string, refreshToken string) (*TokensResponse, error) {
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientSecret)))
	url := "https://login.eveonline.com/oauth/token"
	jsonPayload := []byte(fmt.Sprintf(`{"grant_type":"refresh_token", "refresh_token": "%s"}`, refreshToken))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 8,
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

	var data TokensResponse
	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func getAccessToken(db *sql.DB, character Character, crypt *Crypt, eveClientId string, eveClientSecret string) (string, error) {
	timestamp := getCurrentTimestamp()

	if timestamp > int64(character.expires-1000*60) {
		log.Printf("Updating tokens for Character ID: %d", character.id)

		refreshToken, err := crypt.Decrypt(character.refreshToken)

		if err != nil {
			return "", err
		}

		// TODO: token expired, renew and update
		tokens, err := getTokens(eveClientId, eveClientSecret, refreshToken)

		if err != nil {
			return "", err
		}

		newAccessToken, err := crypt.Encrypt(tokens.AccessToken)

		if err != nil {
			return "", err
		}

		newRefreshToken, err := crypt.Encrypt(tokens.RefreshToken)

		if err != nil {
			return "", err
		}

		_, err = sq.Update("characters").
			Set("accessToken", newAccessToken).
			Set("refreshToken", newRefreshToken).
			Set("expiresAt", tokens.ExpiresIn * 1000 + timestamp).
			Where(sq.Eq{"id": character.id}).
			RunWith(db).
			Exec()

		if err != nil {
			return "", err
		}

		return tokens.AccessToken, nil
	} else {
		token, err := crypt.Decrypt(character.accessToken)

		if err != nil {
			return "", err
		}

		return token, nil
	}
}

func processCharacter(db *sql.DB, character Character, crypt *Crypt) error {
	eveClientId := os.Getenv("EVE_CLIENT_ID")

	if eveClientId == "" {
		return fmt.Errorf("EVE_CLIENT_ID must be set")
	}

	eveClientSecret := os.Getenv("EVE_CLIENT_SECRET")

	if eveClientSecret == "" {
		return fmt.Errorf("EVE_CLIENT_SECRET must be set")
	}

	accessToken, err := getAccessToken(db, character, crypt, eveClientId, eveClientSecret)

	if err != nil {
		return err
	}

	client := NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second * 3)

	pubSub := os.Getenv("PUB_SUB")

	if strings.Contains(character.scopes, "esi-wallet.read_character_wallet.v1") {
		err := processWalletTransactions(client, character.id, pubSub != "")

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func Process() {
	db, err := initializeDb()

	if err != nil {
		log.Fatalf("initializeDb: %v", err)
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("TOKEN_SECRET must be set")
	}

	characters, err := getCharacters(db)

	if err != nil {
		log.Fatal(err)
	}

	crypt := Crypt{key:tokenSecret}

	for _, character := range characters {
		err := processCharacter(db, character, &crypt)
		if err != nil {
			log.Fatalf("Failed to process character ID: %s", character.id)
		}
	}

	defer db.Close()
}

func Esi(ctx context.Context, m PubSubMessage) error {
	Process()
	return nil
}
