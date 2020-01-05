package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"github.com/dariusbakunas/eve-processors/processors"
	"github.com/dariusbakunas/eve-processors/pubsub"
	"log"
	"time"
)

import ps "cloud.google.com/go/pubsub"
import _ "github.com/go-sql-driver/mysql"

func ProcessCharacters() {
	dao, err := db.InitializeDb()

	if err != nil {
		log.Fatalf("initializeDb: %v", err)
	}

	defer dao.Close()

	characters, err := dao.GetCharacters()

	if err != nil {
		log.Fatal(err)
	}

	for _, character := range characters {
		err := processors.ProcessCharacter(dao, character)
		if err != nil {
			log.Fatalf("processors.ProcessCharacter: Failed to process character ID: %d: %v", character.ID, err)
		}
	}
}

func Esi(ctx context.Context, m *ps.Message) error {
	ProcessCharacters()
	return nil
}

type ProcessInit struct {
	dao *db.DB
	esiClient *esi.Client
	characterID int64
}

func initialize(m *ps.Message) (*ProcessInit, error) {
	message := pubsub.Message{}

	if err := json.Unmarshal(m.Data, &message); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	dao, err := db.InitializeDb()

	if err != nil {
		return nil, fmt.Errorf("db.InitializeDb: %v", err)
	}

	accessToken, err := dao.Decrypt(message.AccessToken)

	if err != nil {
		return nil, fmt.Errorf("crypt.Decrypt: %v", err)
	}

	client := esi.NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second * 3)

	return &ProcessInit{
		dao:         dao,
		esiClient:   client,
		characterID: message.CharacterID,
	}, nil
}

func ProcessCharacterWalletTransactions(ctx context.Context, m *ps.Message) error {
	init, err := initialize(m)

	if err != nil {
		return fmt.Errorf("initialize: %v", err)
	}

	defer init.dao.Close()

	err = processors.ProcessWalletTransactions(init.dao, init.esiClient, init.characterID)

	if err != nil {
		return fmt.Errorf("processors.ProcessWalletTransactions: %v", err)
	}

	m.Ack()

	return nil
}

func ProcessCharacterJournalEntries(ctx context.Context, m *ps.Message) error {
	init, err := initialize(m)

	if err != nil {
		return fmt.Errorf("initialize: %v", err)
	}

	defer init.dao.Close()

	err = processors.ProcessJournalEntries(init.dao, init.esiClient, init.characterID)

	if err != nil {
		return fmt.Errorf("processors.ProcessJournalEntries: %v", err)
	}

	m.Ack()

	return nil
}
