package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"gopkg.in/guregu/null.v3"
)

func ProcessBlueprints(dao *db.DB, client *esi.Client, characterID int64) error {
	defer func() {
		dao.CleanupJobLogs("BLUEPRINTS", characterID)
	}()

	blueprintsResponse, err := client.GetBlueprints(characterID, 1)

	if err != nil {
		dao.InsertLogEntry(characterID, "BLUEPRINTS", "FAILURE", "Failed to get blueprints", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetBlueprints: %v", err)
	}

	blueprints := blueprintsResponse.Blueprints

	if blueprintsResponse.Pages > 1 {
		for i := 2; i < blueprintsResponse.Pages; i++ {
			blueprintsResponse, err := client.GetBlueprints(characterID, i)

			if err != nil {
				dao.InsertLogEntry(characterID, "BLUEPRINTS", "FAILURE", "Failed to get blueprints", null.NewString(err.Error(), true))
				return fmt.Errorf("client.GetBlueprints: %v, page: %d", err, i)
			}

			blueprints = append(blueprints, blueprintsResponse.Blueprints...)
		}
	}

	inserted, updated, err := dao.UpdateBlueprints(characterID, blueprints)

	if err != nil {
		dao.InsertLogEntry(characterID, "BLUEPRINTS", "FAILURE", "Failed to update blueprints", null.NewString(err.Error(), true))
		return fmt.Errorf("dao.UpdateBlueprints: %v", err)
	}

	dao.InsertLogEntry(characterID, "BLUEPRINTS", "SUCCESS", fmt.Sprintf("Inserted %d new blueprints and updated %d blueprints", inserted, updated), null.String{})

	return nil
}