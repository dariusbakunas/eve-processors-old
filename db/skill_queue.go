package db

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
)

func (d *DB) UpdateSkillQueue(characterID int64, skillQueueItems []models.SkillQueueItem) (int64, error) {
	if len(skillQueueItems) == 0 {
		return 0, nil
	}

	inserted := int64(0)

	err := d.withTransaction(func(tx *sql.Tx) error {
		_, err := squirrel.Delete("characterSkillQueue").
			Where(squirrel.Eq{"characterId": characterID}).
			RunWith(tx).
			Exec()

		if err != nil {
			return fmt.Errorf("squirrel.Delete: %v", err)
		}

		builder := squirrel.Insert("characterSkillQueue").
			Columns("skillId", "characterId", "finishDate", "finishedLevel", "levelEndSp", "levelStartSp", "queuePosition", "startDate", "trainingStartSp")

		for _, v := range skillQueueItems {
			builder = builder.Values(
				v.SkillID,
				characterID,
				v.FinishDate,
				v.FinishedLevel,
				v.LevelEndSP,
				v.LevelStartSP,
				v.QueuePosition,
				v.StartDate,
				v.TrainingStartSP,
			)
		}

		result, err := builder.RunWith(tx).Exec()

		if err != nil {
			return fmt.Errorf("builder.Exec: %v", err)
		}

		inserted, err = result.RowsAffected()

		if err != nil {
			return fmt.Errorf("RowsAffected: %v", err)
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("d.WithTransaction: %v", err)
	}

	return inserted, nil
}