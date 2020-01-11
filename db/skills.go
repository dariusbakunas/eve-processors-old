package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
)

import sq "github.com/Masterminds/squirrel"

func (d *DB) InsertSkills(characterID int64, skills []models.CharacterSkill) (int, int, error) {
	if len(skills) == 0 {
		return 0, 0, nil
	}

	updated := 0
	inserted := 0

	rows, err := squirrel.Select("skillId").
		From("characterSkills").
		Where(sq.Eq{"characterId": characterID}).
		RunWith(d.db).
		Query()

	if err != nil {
		return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
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
			_, err := squirrel.Insert("characterSkills").
				Columns("skillId", "activeSkillLevel", "skillPointsInSkill", "trainedSkillLevel", "characterId").
				Values(skill.SkillID, skill.ActiveSkillLevel, skill.SP, skill.TrainedLevel, characterID).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Insert: %v", err)
			}

			inserted ++
		}
	}

	return inserted, updated, nil
}