package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
)

import sq "github.com/Masterminds/squirrel"

func (d *DB) InsertSkills(characterID int64, skills []models.CharacterSkill) (int64, int64, error) {
	if len(skills) == 0 {
		return 0, 0, nil
	}

	updated := int64(0)
	inserted := int64(0)

	rows, err := squirrel.Select("skillId").
		From("characterSkills").
		Where(sq.Eq{"characterId": characterID}).
		RunWith(d.db).
		Query()

	if err != nil {
		return 0, 0, fmt.Errorf("squirrel.Select: %v", err)
	}

	defer rows.Close()

	idSet := make(map[int]bool)

	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return 0, 0, fmt.Errorf("rows.Scan: %v", err)
		}

		idSet[id] = true
	}

	builder := squirrel.Insert("characterSkills").
		Columns("skillId", "activeSkillLevel", "skillPointsInSkill", "trainedSkillLevel", "characterId")

	for _, skill := range skills {
		if idSet[skill.SkillID] {
			_, err := squirrel.Update("characterSkills").
				Set("activeSkillLevel", skill.ActiveSkillLevel).
				Set("skillPointsInSkill", skill.SP).
				Set("trainedSkillLevel", skill.TrainedLevel).
				Where(sq.Eq{"characterId": characterID}).
				Where(sq.Eq{"skillId": skill.SkillID}).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
			}

			updated++
		} else {
			builder = builder.Values(skill.SkillID, skill.ActiveSkillLevel, skill.SP, skill.TrainedLevel, characterID)
			inserted ++
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