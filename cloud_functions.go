package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func processWalletTransactions(client *EsiClient, characterId int64) error {
	transactions, err := client.GetWalletTransactions(characterId)

	if err != nil {
		return err
	}

	fmt.Printf("%+v", transactions)

	return nil
}

func Process() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	}

	db, err := initializeDb()

	if err != nil {
		log.Fatalf("initializeDb: %v", err)
	}

	defer db.Close()

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("TOKEN_SECRET must be set")
	}

	characters, err := db.GetCharacters()

	if err != nil {
		log.Fatal(err)
	}

	for _, character := range characters {
		err := ProcessCharacter(db, character, projectID)
		if err != nil {
			log.Fatalf("processCharacter: Failed to process character ID: %d: %v", character.id, err)
		}
	}
}

func Esi(ctx context.Context, m PubSubMessage) error {
	Process()
	return nil
}

func ProcessCharacterWalletTransactions(ctx context.Context, m PubSubMessage) error {
	message := EsiMessage{}

	if err := json.Unmarshal(m.Data, &message); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		return fmt.Errorf("TOKEN_SECRET must be set")
	}

	crypt := Crypt{key:tokenSecret}

	accessToken, err := crypt.Decrypt(message.AccessToken)

	if err != nil {
		return fmt.Errorf("crypt.Decrypt: %v", err)
	}

	client := NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second * 3)

	err = processWalletTransactions(client, message.CharacterID)

	if err != nil {
		return fmt.Errorf("processWalletTransactions: %v", err)
	}

	return nil
}
