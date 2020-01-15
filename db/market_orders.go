package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
	"github.com/dariusbakunas/eve-processors/utils"
)

const BATCH_SIZE = 500

func (d *DB) UpdateMarketOrders(regionID int, orders []models.MarketOrderPriceItem) (int64, int64, error) {
	if len(orders) == 0 {
		return 0, 0, nil
	}

	updated := int64(0)

	rows, err := squirrel.Select("typeId").
		From("marketPrices").
		Join("mapSolarSystems ON mapSolarSystems.solarSystemID = marketPrices.systemId").
		Where(squirrel.Eq{"mapSolarSystems.regionID": regionID}).
		RunWith(d.db).
		Query()

	if err != nil {
		return 0, 0, fmt.Errorf("squirrel.Select: %v", err)
	}

	defer rows.Close()

	idSet, err := d.getIDSet(rows)

	if err != nil {
		return 0, 0, fmt.Errorf("d.getIDSet: %v", err)
	}

	var toInsert []models.MarketOrderPriceItem

	for _, order := range orders {
		if idSet[int64(order.TypeID)] {
			_, err := squirrel.Update("marketPrices").
				Set("buyPrice", order.BuyPrice).
				Set("sellPrice", order.SellPrice).
				Where(squirrel.Eq{"typeId": order.TypeID}).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
			}

			updated++
		} else {
			toInsert = append(toInsert, order)
		}
	}

	for i := 0; i < len(toInsert); i += BATCH_SIZE {
		builder := squirrel.Insert("marketPrices").
			Columns(
				"buyPrice",
				"sellPrice",
				"systemId",
				"typeId",
			)

		batch := toInsert[i:utils.Min(i + BATCH_SIZE, len(toInsert))]

		for _, order := range batch {
			builder = builder.Values(order.BuyPrice, order.SellPrice, order.SystemID, order.TypeID)
		}

		_, err := builder.RunWith(d.db).Exec()

		if err != nil {
			return 0, 0, fmt.Errorf("builder.Exec: %v", err)
		}
	}

	return int64(len(toInsert)), updated, nil
}