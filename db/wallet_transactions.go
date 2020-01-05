package db

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/models"
	"log"
)

import sq "github.com/Masterminds/squirrel"

func (d *DB) InsertWalletTransactions(characterID int64, transactions []models.WalletTransaction) error {
	if len(transactions) == 0 {
		log.Printf("No new transactions for character ID: %d", characterID)
		return nil
	}

	builder := sq.Insert("walletTransactions").
		Options("IGNORE").
		Columns("id", "clientId", "isBuy", "isPersonal", "quantity", "typeId", "locationId", "journalRefId", "unitPrice", "date", "characterId")

	for _, v := range transactions {
		builder = builder.Values(
			v.TransactionId,
			v.ClientId,
			v.IsBuy,
			v.IsPersonal,
			v.Quantity,
			v.TypeId,
			v.LocationId,
			v.JournalRefId,
			v.UnitPrice,
			v.Date,
			characterID)
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
		log.Printf("Inserted %d new transactions for character ID: %d", count, characterID)
	} else {
		log.Printf("No new transactions for character ID: %d", characterID)
	}

	return nil
}