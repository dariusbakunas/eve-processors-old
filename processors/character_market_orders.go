package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"github.com/dariusbakunas/eve-processors/models"
	"gopkg.in/guregu/null.v3"
)

func ProcessCharacterMarketOrders(dao *db.DB, client *esi.Client, characterID int64) error {
	defer func() {
		dao.CleanupJobLogs("MARKET_ORDERS", characterID)
	}()

	activeOrders, err := client.GetMarketOrders(characterID)

	if err != nil {
		dao.InsertLogEntry(characterID, "MARKET_ORDERS", "FAILURE", "Failed to get market orders", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetMarketOrders: %v", err)
	}

	orderHistoryResponse, err := client.GetMarketOrderHistory(characterID, 1)

	if err != nil {
		dao.InsertLogEntry(characterID, "MARKET_ORDERS", "FAILURE", "Failed to get market orders", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetMarketOrderHistory: %v", err)
	}

	orderHistory := orderHistoryResponse.Orders

	if orderHistoryResponse.Pages > 1 {
		for i := 2; i < orderHistoryResponse.Pages; i++ {
			orderHistoryResponse, err := client.GetMarketOrderHistory(characterID, i)

			if err != nil {
				dao.InsertLogEntry(characterID, "MARKET_ORDERS", "FAILURE", "Failed to get market orders", null.NewString(err.Error(), true))
				return fmt.Errorf("client.GetMarketOrderHistory: %v, page: %d", err, i)
			}

			orderHistory = append(orderHistory, orderHistoryResponse.Orders...)
		}
	}

	// orders with duration 0 just replicate market transactions, filter these out
	filteredHistory := make([]models.MarketOrder, 0)
	for _, v := range orderHistory {
		if v.Duration != 0 {
			filteredHistory = append(filteredHistory, v)
		}
	}

	inserted, updated, err := dao.UpdateMarketOrders(characterID, activeOrders, filteredHistory)

	if err != nil {
		dao.InsertLogEntry(characterID, "MARKET_ORDERS", "FAILURE", "Failed to get market orders", null.NewString(err.Error(), true))
		return fmt.Errorf("dao.UpdateMarketOrders: %v", err)
	}

	dao.InsertLogEntry(characterID, "MARKET_ORDERS", "SUCCESS", fmt.Sprintf("Inserted %d new market orders and updated %d orders", inserted, updated), null.String{})

	return nil
}