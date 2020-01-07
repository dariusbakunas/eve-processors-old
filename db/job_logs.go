package db

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"gopkg.in/guregu/null.v3"
	"log"
)

import sq "github.com/Masterminds/squirrel"

func (d *DB) InsertLogEntry(characterID int64, category string, status string, message string, error null.String) {
	_, err := squirrel.Insert("jobLogs").
		Columns( "category", "status", "message", "error", "characterId").
		Values(category, status, message, error, characterID).
		RunWith(d.db).
		Exec()

	if err != nil {
		// log the warning instead of returning the error
		log.Printf("Failed inserting job log entry: %v", err)
	}
}

func (d *DB) CleanupJobLogs(category string, characterID int64) error {
	rows, err := squirrel.
		Select("id").
		From("jobLogs").
		Where(squirrel.Eq{"category": category}).
		OrderBy("createdAt DESC").Limit(3).
		RunWith(d.db).
		Query()

	if err != nil {
		return fmt.Errorf("squirrel.Select: %v", err)
	}

	var ids []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return fmt.Errorf("rows.Scan: %v", err)
		}

		ids = append(ids, id)
	}

	_, err = squirrel.Delete("jobLogs").
		Where(sq.NotEq{"id": ids}).
		Where(sq.Eq{"category": category}).
		Where(sq.Eq{"characterId": characterID}).
		RunWith(d.db).
		Exec()

	if err != nil {
		return fmt.Errorf("squirrel.Delete: %v", err)
	}

	return nil
}