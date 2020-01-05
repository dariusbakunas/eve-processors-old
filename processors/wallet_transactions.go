package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
)

func ProcessWalletTransactions(dao *db.DB, client *esi.Client, characterId int64) error {
	transactions, err := client.GetWalletTransactions(characterId)

	if err != nil {
		return fmt.Errorf("client.GetWalletTransactions: %v", err)
	}

	err = dao.InsertWalletTransactions(characterId, transactions)

	if err != nil {
		return fmt.Errorf("dao.InsertWalletTransactions: %v", err)
	}

	return nil
}