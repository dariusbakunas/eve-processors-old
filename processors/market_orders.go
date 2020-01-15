package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"github.com/dariusbakunas/eve-processors/models"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"time"
)

const FORGE_REGION_ID = 10000002
const JITA_SYSTEM_ID = 30000142

func ConvertOrderMap(orders map[int][]models.MarketOrder) []models.MarketOrderPriceItem {
	var result []models.MarketOrderPriceItem

	for typeID, typeOrders := range orders {
		maxBuyOrder := models.MarketOrder{
			Duration:     0,
			IsBuy:        true,
			Issued:       time.Time{},
			LocationID:   0,
			MinVolume:    0,
			OrderID:      0,
			Price:        decimal.NewFromInt(0),
			SystemID:     JITA_SYSTEM_ID,
			Range:        "",
			TypeID:       0,
			VolumeRemain: 0,
			VolumeTotal:  0,
		}
		minSellOrder := models.MarketOrder{
			Duration:     0,
			IsBuy:        false,
			Issued:       time.Time{},
			LocationID:   0,
			MinVolume:    0,
			OrderID:      0,
			Price:        decimal.NewFromInt32(math.MaxInt32),
			SystemID:     JITA_SYSTEM_ID,
			Range:        "",
			TypeID:       0,
			VolumeRemain: 0,
			VolumeTotal:  0,
		}
		for _, order := range typeOrders {
			if order.SystemID != JITA_SYSTEM_ID {
				continue
			}

			if order.IsBuy {
				if maxBuyOrder.Price.LessThan(order.Price) {
					maxBuyOrder = order
				}
			} else {
				if minSellOrder.Price.GreaterThan(order.Price) {
					minSellOrder = order
				}
			}
		}

		if maxBuyOrder.SystemID == JITA_SYSTEM_ID && minSellOrder.SystemID == JITA_SYSTEM_ID {
			result = append(result, models.MarketOrderPriceItem{
				TypeID:    typeID,
				BuyPrice:  minSellOrder.Price,
				SellPrice: maxBuyOrder.Price,
				SystemID:  JITA_SYSTEM_ID,
			})
		}
	}

	return result
}

func ProcessMarketOrders(dao *db.DB, client *esi.Client) error {
	marketOrders, err := client.GetMarketOrders(FORGE_REGION_ID)

	if err != nil {
		return fmt.Errorf("client.GetMarketOrders: %v", err)
	}

	priceItems := ConvertOrderMap(marketOrders)

	inserted, updated, err := dao.UpdateMarketOrders(FORGE_REGION_ID, priceItems)

	if err != nil {
		return fmt.Errorf("dao.UpdateMarketOrders: %v", err)
	}

	log.Printf("Inserted %d new market orders, updated: %d", inserted, updated)

	return nil
}