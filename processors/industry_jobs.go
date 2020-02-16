package processors

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/db"
	"github.com/dariusbakunas/eve-processors/esi"
	"gopkg.in/guregu/null.v3"
)

func ProcessCharacterIndustryJobs(dao *db.DB, client *esi.Client, characterID int64) error {
	defer func() {
		dao.CleanupJobLogs("INDUSTRY_JOBS", characterID)
	}()

	industryJobs, err := client.GetIndustryJobs(characterID)

	if err != nil {
		dao.InsertLogEntry(characterID, "INDUSTRY_JOBS", "FAILURE", "Failed to get industry jobs", null.NewString(err.Error(), true))
		return fmt.Errorf("client.GetIndustryJobs: %v", err)
	}

	inserted, updated, err := dao.UpdateIndustryJobs(characterID, industryJobs)

	if err != nil {
		dao.InsertLogEntry(characterID, "INDUSTRY_JOBS", "FAILURE", "Failed to update industry jobs", null.NewString(err.Error(), true))
		return fmt.Errorf("dao.UpdateIndustryJobs: %v", err)
	}

	dao.InsertLogEntry(characterID, "INDUSTRY_JOBS", "SUCCESS", fmt.Sprintf("Inserted %d new industry jobs and updated %d industry jobs", inserted, updated), null.String{})
	return nil
}