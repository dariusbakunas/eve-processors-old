package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
)

func (d *DB) UpdateBlueprints(characterID int64, blueprints []models.Blueprint) (int64, int64, error) {
	if len(blueprints) == 0 {
		return 0, 0, nil
	}

	rows, err := squirrel.Select("id").
		From("blueprints").
		Where(squirrel.Eq{"characterId": characterID}).
		RunWith(d.db).
		Query()

	if err != nil {
		return 0, 0, fmt.Errorf("squirrel.Select: %v", err)
	}

	defer rows.Close()

	idSet := make(map[int64]bool)

	for rows.Next() {
		var id int64
		err := rows.Scan(&id)

		if err != nil {
			return 0, 0, fmt.Errorf("rows.Scan: %v", err)
		}

		idSet[id] = true
	}

	updated := int64(0)
	inserted := int64(0)

	builder := squirrel.Insert("blueprints").Columns(
		"id",
		"locationType",
		"locationId",
		"materialEfficiency",
		"isCopy",
		"maxRuns",
		"timeEfficiency",
		"typeId",
		"characterId",
		)

	for _, blueprint := range blueprints {
		if idSet[blueprint.ID] {
			_, err := squirrel.Update("blueprints").
				Set("locationType", blueprint.LocationFlag).
				Set("locationId", blueprint.LocationID).
				Set("materialEfficiency", blueprint.ME).
				Set("maxRuns", blueprint.Runs).
				Set("timeEfficiency", blueprint.TE).
				Where(squirrel.Eq{"id": blueprint.ID}).
				Where(squirrel.Eq{"characterId": characterID}).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
			}

			updated++
		} else {
			isCopy := false

			if blueprint.Runs != -1 {
				isCopy = true
			}

			builder = builder.Values(
				blueprint.ID,
				blueprint.LocationFlag,
				blueprint.LocationID,
				blueprint.ME,
				isCopy,
				blueprint.Runs,
				blueprint.TE,
				blueprint.TypeID,
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

	return inserted, updated, nil
}