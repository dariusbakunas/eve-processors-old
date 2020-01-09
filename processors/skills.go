package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"gopkg.in/guregu/null.v3"
)

func ProcessSkills(dao *db.DB, client *esi.Client, characterID int64) error {
	defer func() {
		dao.CleanupJobLogs("SKILLS", characterID)
	}()

	skillsResponse, err := client.GetSkills(characterID)

	if err != nil {
		dao.InsertLogEntry(characterID, "SKILLS", "FAILURE", "Failed to get skills", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetSkills: %v", err)
	}

	err = dao.UpdateCharacterSP(characterID, skillsResponse.TotalSP)

	if err != nil {
		dao.InsertLogEntry(characterID, "SKILLS", "FAILURE", "Failed to update character SP", null.NewString(err.Error(), true))
		return fmt.Errorf("dao.UpdateCharacterSP: %v", err)
	}

	inserted, updated, err := dao.InsertSkills(characterID, skillsResponse.Skills)

	if err != nil {
		dao.InsertLogEntry(characterID, "SKILLS", "FAILURE", "Failed to insert skills", null.NewString(err.Error(), true))
	}

	dao.InsertLogEntry(characterID, "SKILLS", "SUCCESS", fmt.Sprintf("Inserted %d new skills and updated %d skills", inserted, updated), null.String{})

	return nil
}