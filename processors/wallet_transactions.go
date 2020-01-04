package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/esi"
)

func ProcessWalletTransactions(client *esi.Client, characterId int64) error {
	transactions, err := client.GetWalletTransactions(characterId)

	if err != nil {
		return err
	}

	fmt.Printf("%+v", transactions)

	return nil
}