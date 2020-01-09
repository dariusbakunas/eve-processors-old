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

	for _, skill := range skills {
		result, err := squirrel.Update("characterSkills").
			Set("activeSkillLevel", skill.ActiveSkillLevel).
			Set("skillPointsInSkill", skill.SP).
			Set("trainedSkillLevel", skill.SP).
			Where(sq.Eq{"characterId": characterID}).
			Where(sq.Eq{"skillId": skill.SkillID}).
			RunWith(d.db).
			Exec()

		if err != nil {
			return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
		}

		count, err := result.RowsAffected()

		if err != nil {
			return 0, 0, fmt.Errorf("result.RowsAffected: %v", err)
		}

		if count == 0 {
			_, err := squirrel.Insert("characterSkills").
				Columns("skillId", "activeSkillLevel", "skillPointsInSkill", "trainedSkillLevel", "characterId").
				Values(skill.SkillID, skill.ActiveSkillLevel, skill.SP, skill.TrainedLevel, characterID).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Insert: %v", err)
			}
		}
	}

	return inserted, updated, nil
}