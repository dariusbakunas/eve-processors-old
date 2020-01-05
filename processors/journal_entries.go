package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
)

func ProcessJournalEntries(dao *db.DB, client *esi.Client, characterId int64) error {
	journalEntriesResponse, err := client.GetJournalEntries(characterId, 1)

	if err != nil {
		return fmt.Errorf("client.GetJournalEntries: %v", err)
	}

	transactions := journalEntriesResponse.Entries

	if journalEntriesResponse.Pages > 1 {
		for i := 2; i < journalEntriesResponse.Pages; i++ {
			journalEntriesResponse, err := client.GetJournalEntries(characterId, i)

			if err != nil {
				return fmt.Errorf("client.GetJournalEntries: %v, page: %d", err, i)
			}

			transactions = append(transactions, journalEntriesResponse.Entries...)
		}
	}

	err = dao.InsertJournalEntries(characterId, transactions)

	if err != nil {
		return fmt.Errorf("dao.InsertJournalEntries: %v", err)
	}

	return nil
}