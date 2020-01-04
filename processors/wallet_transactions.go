package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
)

func ProcessWalletTransactions(dao *db.DB, client *esi.Client, characterId int64) error {
	transactions, err := client.GetWalletTransactions(characterId)

	if err != nil {
		return err
	}

	t := make([]db.WalletTransaction, len(transactions))

	for i, v := range transactions {
		t[i] = db.WalletTransaction{
			ID:           v.TransactionId,
			ClientID:     v.ClientId,
			CharacterID:  characterId,
			IsBuy:        v.IsBuy,
			IsPersonal:   v.IsPersonal,
			Quantity:     v.Quantity,
			TypeID:       v.TypeId,
			LocationID:   v.LocationId,
			JournalRefID: v.JournalRefId,
			UnitPrice:    v.UnitPrice,
			Date:         v.Date,
		}
	}

	err = dao.InsertWalletTransactions(characterId, t)

	if err != nil {
		return fmt.Errorf("dao.InsertWalletTransactions: %v", err)
	}

	return nil
}