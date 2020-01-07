package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"gopkg.in/guregu/null.v3"
	"log"
)

func ProcessWalletTransactions(dao *db.DB, client *esi.Client, characterID int64) error {
	defer func() {
		err := dao.CleanupJobLogs("WALLET_TRANSACTIONS", characterID)

		if err != nil {
			log.Printf("d.CleanupJobLogs: %v", err)
		}
	}()

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
		log.Printf("Inserted %d new wallet transactions for character ID: %d", count, characterID)
		dao.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "SUCCESS", fmt.Sprintf("Inserted %d new transactions", count), null.String{})
	} else {
		log.Printf("No new wallet transactions for character ID: %d", characterID)
		dao.InsertLogEntry(characterID, "WALLET_TRANSACTIONS", "SUCCESS", "No new transactions", null.String{})
	}

	return nil
}