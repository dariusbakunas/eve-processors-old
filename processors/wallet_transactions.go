package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"gopkg.in/guregu/null.v3"
)

func ProcessWalletTransactions(dao *db.DB, client *esi.Client, characterID int64) error {
	transactions, err := client.GetWalletTransactions(characterID)

	if err != nil {
		dao.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "FAILURE", "Failed to get wallet transactions", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetWalletTransactions: %v", err)
	}

	count, err := dao.InsertWalletTransactions(characterID, transactions)

	if err != nil {
		dao.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "FAILURE", "Failed to get wallet transactions", null.NewString(err.Error(), true))
		return fmt.Errorf("dao.InsertWalletTransactions: %v", err)
	}

	if count > 0 {
		dao.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "SUCCESS", fmt.Sprintf("Inserted %d new transactions", count), null.String{})
	} else {

	}

	return nil
}