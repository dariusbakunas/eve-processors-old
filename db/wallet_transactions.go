package db

import (
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

import sq "github.com/Masterminds/squirrel"

type WalletTransaction struct {
	ID int64
	CharacterID int64
	ClientID int64
	IsBuy bool
	IsPersonal bool
	Quantity int64
	TypeID int
	LocationID int64
	JournalRefID int64
	UnitPrice decimal.Decimal
	Date time.Time
}

func (d *DB) InsertWalletTransactions(characterID int64, transactions []WalletTransaction) error {
	if len(transactions) == 0 {
		log.Printf("No new transactions for character ID: %d", characterID)
		return nil
	}

	builder := sq.Insert("walletTransactions").
		Options("IGNORE").
		Columns("id", "clientId", "isBuy", "isPersonal", "quantity", "typeId", "locationId", "journalRefId", "unitPrice", "date", "characterId")

	for _, v := range transactions {
		builder = builder.Values(
			v.ID,
			v.ClientID,
			v.IsBuy,
			v.IsPersonal,
			v.Quantity,
			v.TypeID,
			v.LocationID,
			v.JournalRefID,
			v.UnitPrice,
			v.Date,
			v.CharacterID)
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