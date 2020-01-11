package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"gopkg.in/guregu/null.v3"
)

func ProcessSkillQueue(dao *db.DB, client *esi.Client, characterID int64) error {
	defer func() {
		dao.CleanupJobLogs("SKILL_QUEUE", characterID)
	}()

	skillQueueItems, err := client.GetSkillQueue(characterID)

	if err != nil {
		dao.InsertLogEntry(characterID, "SKILL_QUEUE", "FAILURE", "Failed to get skill queue", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetSkillQueue: %v", err)
	}

	inserted, err := dao.UpdateSkillQueue(characterID, skillQueueItems)

	if err != nil {
		dao.InsertLogEntry(characterID, "SKILL_QUEUE", "FAILURE", "Failed to update skill queue", null.NewString(err.Error(), true))
		return fmt.Errorf("dao.UpdateSkillQueue: %v", err)
	}

	dao.InsertLogEntry(characterID, "SKILL_QUEUE", "SUCCESS", fmt.Sprintf("Inserted %d skill queue items", inserted), null.String{})
	return nil
}