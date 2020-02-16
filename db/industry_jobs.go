package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
	"time"
)

type activeJob struct {
	ID int64
	Duration int
	Started time.Time
}

func (d *DB) getActiveJobs(characterID int64) (map[int64]activeJob, error) {
	rows, err := squirrel.Select("id", "duration", "startDate").
		From("industryJobs").
		Where(squirrel.Eq{"installerId": characterID}).
		Where(squirrel.Eq{"status": "active"}).
		RunWith(d.db).
		Query()

	if err != nil {
		return nil, fmt.Errorf("squirrel.Select: %v", err)
	}

	defer rows.Close()

	activeJobs := make(map[int64]activeJob)

	for rows.Next() {
		var job activeJob

		err := rows.Scan(&job.ID, &job.Duration, &job.Started)

		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %v", err)
		}

		activeJobs[job.ID] = job
	}

	return activeJobs, nil
}

func (d *DB) UpdateIndustryJobs(characterID int64, jobs []models.IndustryJob) (int64, int64, error) {
	activeJobs, err := d.getActiveJobs(characterID)

	if err != nil {
		return 0, 0, fmt.Errorf("d.getActiveJobs: %v", err)
	}

	builder := squirrel.Insert("industryJobs").
		Options("IGNORE").
		Columns(
			"id",
			"activityId",
			"blueprintId",
			"blueprintLocationId",
			"blueprintTypeId",
			"completedCharacterId",
			"completedDate",
			"cost",
			"duration",
			"endDate",
			"facilityId",
			"installerId",
			"licensedRuns",
			"outputLocationId",
			"pauseDate",
			"probability",
			"productTypeId",
			"runs",
			"startDate",
			"stationId",
			"status",
			"successfulRuns")

	updated := int64(0)
	inserted := int64(0)

	for _, i := range jobs {
		if _, ok := activeJobs[i.JobID]; ok {
			_, err = squirrel.Update("industryJobs").
				Set("completedCharacterId", i.CompletedCharacterID).
				Set("completedDate", i.CompletedDate).
				Set("endDate", i.EndDate).
				Set("pauseDate", i.PauseDate).
				Set("status", i.Status).
				Set("successfulRuns", i.SuccessfulRuns).
				Where(squirrel.Eq{"id": i.JobID}).
				RunWith(d.db).
				Exec()

			if err != nil {
				return 0, 0, fmt.Errorf("squirrel.Update: %v", err)
			}

			updated++
			delete(activeJobs, i.JobID)
		} else {
			builder = builder.Values(
				i.JobID,
				i.ActivityID,
				i.BlueprintID,
				i.BlueprintLocationID,
				i.BlueprintTypeID,
				i.CompletedCharacterID,
				i.CompletedDate,
				i.Cost,
				i.Duration,
				i.EndDate,
				i.FacilityID,
				i.InstalledID,
				i.LicensedRuns,
				i.OutputLocationID,
				i.PauseDate,
				i.Probability,
				i.ProductTypeID,
				i.Runs,
				i.StartDate,
				i.StationID,
				i.Status,
				i.SuccessfulRuns,
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