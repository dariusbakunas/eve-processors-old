package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
	"log"
)

func (d *DB) InsertJournalEntries(characterID int64, entries []models.JournalEntry) error {
	if len(entries) == 0 {
		log.Printf("No new journal entries for character ID: %d", characterID)
		return nil
	}

	builder := squirrel.Insert("journalEntries").
		Options("IGNORE").
		Columns("id", "amount", "balance", "contextId", "contextIdType", "date", "description", "firstPartyId", "reason", "refType", "secondPartyId", "tax", "taxReceiverId", "characterId")

	for _, v := range entries {
		builder = builder.Values(
			v.ID,
			v.Amount,
			v.Balance,
			v.ContextID,
			v.ContextIDType,
			v.Date,
			v.Description,
			v.FirstPartyID,
			v.Reason,
			v.RefType,
			v.SecondPartyID,
			v.Tax,
			v.TaxReceiverID,
			characterID,
			)
	}

	result, err := builder.RunWith(d.db).Exec()

	if err != nil {
		return fmt.Errorf("builder.Exec: %v", err)
	}

	count, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("result.RowsAffected: %v", err)
	}

	if count > 0 {
		log.Printf("Inserted %d new journal entries for character ID: %d", count, characterID)
	} else {
		log.Printf("No new journal entries for character ID: %d", characterID)
	}

	return nil
}