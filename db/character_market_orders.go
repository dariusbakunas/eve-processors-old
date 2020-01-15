package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
	"gopkg.in/guregu/null.v3"
	"time"
)

type activeOrder struct {
	ID int64
	Duration int
	Issued time.Time
}

func (d *DB) getActiveOrders(characterID int64) (map[int64]activeOrder, error) {
	rows, err := squirrel.Select("id", "duration", "issued").
		From("characterMarketOrders").
		Where(squirrel.Eq{"characterId": characterID}).
		Where(squirrel.Eq{"state": "active"}).
		RunWith(d.db).
		Query()

	if err != nil {
		return nil, fmt.Errorf("squirrel.Select: %v", err)
	}

	defer rows.Close()

	activeOrders := make(map[int64]activeOrder)

	for rows.Next() {
		var order activeOrder

		err := rows.Scan(&order.ID, &order.Duration, &order.Issued)

		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %v", err)
		}

		activeOrders[order.ID] = order
	}

	return activeOrders, nil
}

func (d *DB) UpdateCharacterMarketOrders(characterID int64, orders []models.CharacterMarketOrder, history []models.CharacterMarketOrder) (int64, int64, error) {
	allOrders := append(orders, history...)

	activeOrders, err := d.getActiveOrders(characterID)

	builder := squirrel.Insert("characterMarketOrders").
		Options("IGNORE").
		Columns(
			"id",
			"duration",
			"escrow",
			"isBuy",
			"isCorporation",
			"issued",
			"locationId",
			"minVolume",
			"price",
			"range",
			"regionId",
			"typeId",
			"state",
			"volumeRemain",
			"volumeTotal",
			"characterId",
			)

	updated := int64(0)
	inserted := int64(0)

	for _, o := range allOrders {
		state := o.State

		if !o.State.Valid {
			state = null.NewString("active", true)
		}

		if _, ok := activeOrders[o.OrderID]; ok {
			_, err = squirrel.Update("characterMarketOrders").
				Set("price", o.Price).
				Set("state", state).
				Set("volumeRemain", o.VolumeRemain).
				Where(squirrel.Eq{"id": o.OrderID}).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
			}

			updated++

			delete(activeOrders, o.OrderID)
		} else if o.Duration == 0 {
			builder = builder.Values(
				o.OrderID,
				o.Duration,
				o.Escrow,
				o.IsBuy,
				o.IsCorporation,
				o.Issued,
				o.LocationID,
				o.MinVolume,
				o.Price,
				o.Range,
				o.RegionID,
				o.TypeID,
				state,
				o.VolumeRemain,
				o.VolumeTotal,
				characterID,
				)

			inserted++
		}
	}

	if inserted > 0 {
		result, err := builder.RunWith(d.db).Exec()

		if err != nil {
			return 0, 0, fmt.Errorf("builder.Exec: %v", err)
		}

		inserted, err = result.RowsAffected()

		if err != nil {
			return 0, 0, fmt.Errorf("result.RowsAffected: %v", err)
		}
	}

	// TODO: expire any active orders
	for _, order := range activeOrders {
		expires := order.Issued.AddDate(0, 0, order.Duration)

		if expires.Before(time.Now()) {
			_, err = squirrel.Update("characterMarketOrders").
				Set("state", "expired"). // no way to know if this is expired or cancelled, since it was not part of history
				Where("id", order.ID).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
			}

			updated++
		}
	}

	return inserted, updated, nil
}