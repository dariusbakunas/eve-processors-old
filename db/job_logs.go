package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"gopkg.in/guregu/null.v3"
)

func (d *DB) InsertLogEntry(characterID int64, category string, status string, message string, error null.String) error {
	_, err := squirrel.Insert("jobLogs").
		Columns( "category", "status", "message", "error", "characterId").
		Values(category, status, message, error, characterID).
		RunWith(d.db).
		Exec()

	if err != nil {
		return fmt.Errorf("builder.Exec: %v", err)
	}

	return nil
}