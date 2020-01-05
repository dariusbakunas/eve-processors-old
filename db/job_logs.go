package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/dariusbakunas/eve-processors/models"
)

func (d *DB) InsertLogEntry(characterID int64, entry models.JobLogEntry) error {
	_, err := squirrel.Insert("jobLogs").
		Columns( "status", "message", "error", "characterId").
		Values(entry.Status, entry.Message, entry.Error, entry.CharacterID).
		RunWith(d.db).
		Exec()

	if err != nil {
		return fmt.Errorf("builder.Exec: %v", err)
	}

	return nil
}