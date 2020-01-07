package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"gopkg.in/guregu/null.v3"
	"log"
)

func ProcessJournalEntries(dao *db.DB, client *esi.Client, characterID int64) error {
	defer func() {
		err := dao.CleanupJobLogs("WALLET_JOURNAL", characterID)

		if err != nil {
			log.Printf("d.CleanupJobLogs: %v", err)
		}
	}()

	journalEntriesResponse, err := client.GetJournalEntries(characterID, 1)

	if err != nil {
		dao.InsertLogEntry(characterID, "WALLET_JOURNAL", "FAILURE", "Failed to get journal entries", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetJournalEntries: %v", err)
	}

	transactions := journalEntriesResponse.Entries

	if journalEntriesResponse.Pages > 1 {
		for i := 2; i < journalEntriesResponse.Pages; i++ {
			journalEntriesResponse, err := client.GetJournalEntries(characterID, i)

			if err != nil {
				dao.InsertLogEntry(characterID, "WALLET_JOURNAL", "FAILURE", "Failed to get journal entries", null.NewString(err.Error(), true))
				return fmt.Errorf("client.GetJournalEntries: %v, page: %d", err, i)
			}

			transactions = append(transactions, journalEntriesResponse.Entries...)
		}
	}

	count, err := dao.InsertJournalEntries(characterID, transactions)

	if err != nil {
		dao.InsertLogEntry(characterID, "WALLET_JOURNAL", "FAILURE", "Failed to get journal entries", null.NewString(err.Error(), true))
		return fmt.Errorf("dao.InsertJournalEntries: %v", err)
	}

	if count > 0 {
		log.Printf("Inserted %d new journal entries for character ID: %d", count, characterID)
		dao.InsertLogEntry(characterID, "WALLET_JOURNAL", "SUCCESS", fmt.Sprintf("Inserted %d new journal entries", count), null.String{})
	} else {
		log.Printf("No new journal entries for character ID: %d", characterID)
		dao.InsertLogEntry(characterID, "WALLET_JOURNAL", "SUCCESS", "No new journal entries", null.String{})
	}

	return nil
}