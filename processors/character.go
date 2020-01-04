package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"github.com/dariusbakunas/eve-processors/pubsub"
	"os"
	"strings"
	"time"
)

func ProcessCharacter(db *db.DB, character db.Character, projectID string) error {
	eveClientId := os.Getenv("EVE_CLIENT_ID")

	if eveClientId == "" {
		return fmt.Errorf("EVE_CLIENT_ID must be set")
	}

	eveClientSecret := os.Getenv("EVE_CLIENT_SECRET")

	if eveClientSecret == "" {
		return fmt.Errorf("EVE_CLIENT_SECRET must be set")
	}

	accessToken, err := esi.GetAccessToken(db, character, eveClientId, eveClientSecret)

	if err != nil {
		return fmt.Errorf("getAccessToken: %v", err);
	}

	client := esi.NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second * 3)

	if strings.Contains(character.Scopes, "esi-wallet.read_character_wallet.v1") {
		topicID := os.Getenv("PUBSUB_WALLET_TRANSACTIONS_TOPIC_ID")

		encryptedToken, err := db.Encrypt(accessToken)

		if err != nil {
			return fmt.Errorf("db.Encrypt: %v", err)
		}

		if topicID != "" {
			err := pubsub.PublishMessage(projectID, topicID, character.ID, encryptedToken)

			if err != nil {
				return fmt.Errorf("PublishMessage: %v", err)
			}
		} else {
			err := ProcessWalletTransactions(client, character.ID)

			if err != nil {
				return fmt.Errorf("processWalletTransactions: %v", err)
			}
		}
	}

	return nil
}