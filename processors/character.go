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

type FnRef struct {
	fn   func(dao *db.DB, client *esi.Client, characterID int64) error
	name string
}

func ProcessCharacter(dao *db.DB, character db.Character) error {
	messageMap := map[string][]string{
		"esi-wallet.read_character_wallet.v1":  {"PUBSUB_WALLET_TRANSACTIONS_TOPIC_ID", "PUBSUB_JOURNAL_ENTRIES_TOPIC_ID"},
		"esi-skills.read_skills.v1":            {"PUBSUB_SKILLS_ID", "PUBSUB_SKILL_QUEUE_ID"},
		"esi-markets.read_character_orders.v1": {"PUBSUB_CHARACTER_MARKET_ORDERS_ID"},
		"esi-characters.read_blueprints.v1":    {"PUBSUB_CHARACTER_BLUEPRINTS_ID"},
	}

	fnMap := map[string][]FnRef {
		"esi-wallet.read_character_wallet.v1":  {FnRef{fn: ProcessWalletTransactions, name: "ProcessWalletTransactions" }, FnRef{fn: ProcessJournalEntries, name: "ProcessJournalEntries" }},
		"esi-skills.read_skills.v1":            {FnRef{fn: ProcessSkills, name: "ProcessSkills"}, FnRef{fn: ProcessSkillQueue, name: "ProcessSkillQueue"}},
		"esi-markets.read_character_orders.v1": {FnRef{fn: ProcessCharacterMarketOrders, name: "ProcessCharacterMarketOrders"}},
		"esi-characters.read_blueprints.v1":    {FnRef{fn: ProcessBlueprints, name: "ProcessBlueprints"}},
	}

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

	client := esi.NewEsiClient("https://esi.evetech.net/latest", accessToken, time.Second*3)
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	scopes := strings.Split(character.Scopes, " ")

	for _, scope := range scopes {
		if projectID != "" {
			if messages, ok := messageMap[scope]; ok {
				for _, message := range messages {
					err = pubsub.PublishMessage(dao, projectID, message, character.ID, accessToken)

					if err != nil {
						return fmt.Errorf("PublishMessage %s: %v", message, err)
					}
				}
			}
		} else {
			if functions, ok := fnMap[scope]; ok {
				for _, fnRef := range functions {
					err = fnRef.fn(dao, client, character.ID)

					if err != nil {
						return fmt.Errorf("%s: %v", fnRef.name, err)
					}
				}
			}
		}
	}

	return nil
}
