package db

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/models"
	"gopkg.in/guregu/null.v3"
	"log"
)

import sq "github.com/Masterminds/squirrel"

func (d *DB) InsertWalletTransactions(characterID int64, transactions []models.WalletTransaction) error {
	defer func() {
		err := d.Cleanup("WALLET_TRANSACTIONS", characterID)

		if err != nil {
			log.Printf("d.Cleanup: %v", err)
		}
	}()

	if len(transactions) == 0 {
		d.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "SUCCESS", "No new transactions", null.String{})
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
		d.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "FAILURE", "Failed to get wallet transactions", null.NewString(err.Error(), true))
		return fmt.Errorf("builder.Exec: %v", err)
	}

	count, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("result.RowsAffected: %v", err)
	}

	if count > 0 {
		d.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "SUCCESS", fmt.Sprintf("Inserted %d new transactions", count), null.String{})
		log.Printf("Inserted %d new transactions for character ID: %d", count, characterID)
	} else {
		d.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "SUCCESS", "No new transactions", null.String{})
		log.Printf("No new transactions for character ID: %d", characterID)
	}

	return nil
}