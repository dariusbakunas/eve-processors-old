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

func ProcessCharacter(dao *db.DB, character db.Character) error {
	eveClientId := os.Getenv("EVE_CLIENT_ID")

	if eveClientId == "" {
		return fmt.Errorf("EVE_CLIENT_ID must be set")
	}

	eveClientSecret := os.Getenv("EVE_CLIENT_SECRET")

	if eveClientSecret == "" {
		return fmt.Errorf("EVE_CLIENT_SECRET must be set")
	}

	accessToken, err := esi.GetAccessToken(dao, character, eveClientId, eveClientSecret)

	if err != nil {
		return fmt.Errorf("getAccessToken: %v", err);
	}

	client := esi.NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second * 3)
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if strings.Contains(character.Scopes, "esi-wallet.read_character_wallet.v1") {
		if projectID != "" {
			err = pubsub.PublishMessage(dao, projectID, "PUBSUB_WALLET_TRANSACTIONS_TOPIC_ID", character.ID, accessToken)

			if err != nil {
				return fmt.Errorf("PublishMessage PUBSUB_WALLET_TRANSACTIONS_TOPIC_ID: %v", err)
			}

			err = pubsub.PublishMessage(dao, projectID, "PUBSUB_JOURNAL_ENTRIES_TOPIC_ID", character.ID, accessToken)

			if err != nil {
				return fmt.Errorf("PublishMessage PUBSUB_JOURNAL_ENTRIES_TOPIC_ID: %v", err)
			}
		} else {
			err := ProcessWalletTransactions(dao, client, character.ID)

			if err != nil {
				return fmt.Errorf("ProcessWalletTransactions: %v", err)
			}

			err = ProcessJournalEntries(dao, client, character.ID)

			if err != nil {
				return fmt.Errorf("ProcessJournalEntries: %v", err)
			}
		}
	}

	if strings.Contains(character.Scopes, "esi-skills.read_skills.v1") {
		if projectID != "" {
			err = pubsub.PublishMessage(dao, projectID, "PUBSUB_SKILLS_ID", character.ID, accessToken)

			if err != nil {
				return fmt.Errorf("PublishMessage PUBSUB_SKILLS_ID: %v", err)
			}

			err = pubsub.PublishMessage(dao, projectID, "PUBSUB_SKILL_QUEUE_ID", character.ID, accessToken)

			if err != nil {
				return fmt.Errorf("PublishMessage PUBSUB_SKILL_QUEUE_ID: %v", err)
			}
		} else {
			err = ProcessSkills(dao, client, character.ID)

			if err != nil {
				return fmt.Errorf("ProcessSkills: %v", err)
			}

			err = ProcessSkillQueue(dao, client, character.ID)

			if err != nil {
				return fmt.Errorf("ProcessSkillQueue: %v", err)
			}
		}
	}

	return nil
}