package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"github.com/dariusbakunas/eve-processors/processors"
	"github.com/dariusbakunas/eve-processors/pubsub"
	"github.com/dariusbakunas/eve-processors/utils"
	"log"
	"os"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func ProcessCharacters() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	}

	db, err := db.InitializeDb()

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
		err := processors.ProcessCharacter(db, character, projectID)
		if err != nil {
			log.Fatalf("processCharacter: Failed to process character ID: %d: %v", character.ID, err)
		}
	}
}

func Esi(ctx context.Context, m PubSubMessage) error {
	ProcessCharacters()
	return nil
}

func ProcessCharacterWalletTransactions(ctx context.Context, m PubSubMessage) error {
	message := pubsub.Message{}

	if err := json.Unmarshal(m.Data, &message); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		return fmt.Errorf("TOKEN_SECRET must be set")
	}

	crypt := utils.Crypt{Key: tokenSecret}

	accessToken, err := crypt.Decrypt(message.AccessToken)

	if err != nil {
		return fmt.Errorf("crypt.Decrypt: %v", err)
	}

	client := esi.NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second * 3)

	err = processors.ProcessWalletTransactions(client, message.CharacterID)

	if err != nil {
		return fmt.Errorf("processWalletTransactions: %v", err)
	}

	return nil
}
