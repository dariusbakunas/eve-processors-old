package esi

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func ProcessCharacter(db *DB, character Character, projectID string) error {
	eveClientId := os.Getenv("EVE_CLIENT_ID")

	if eveClientId == "" {
		return fmt.Errorf("EVE_CLIENT_ID must be set")
	}

	eveClientSecret := os.Getenv("EVE_CLIENT_SECRET")

	if eveClientSecret == "" {
		return fmt.Errorf("EVE_CLIENT_SECRET must be set")
	}

	accessToken, err := GetAccessToken(db, character, eveClientId, eveClientSecret)

	if err != nil {
		return fmt.Errorf("getAccessToken: %v", err);
	}

	client := NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second * 3)

	if strings.Contains(character.scopes, "esi-wallet.read_character_wallet.v1") {
		topicID := os.Getenv("PUBSUB_WALLET_TRANSACTIONS_TOPIC_ID")

		encryptedToken, err := db.crypt.Encrypt(accessToken)

		if err != nil {
			return fmt.Errorf("crypt.Encrypt: %v", err)
		}

		if topicID != "" {
			err := PublishMessage(projectID, topicID, character.id, encryptedToken)

			if err != nil {
				return fmt.Errorf("PublishMessage: %v", err)
			}
		} else {
			err := processWalletTransactions(client, character.id)

			if err != nil {
				return fmt.Errorf("processWalletTransactions: %v", err)
			}
		}
	}

	return nil
}