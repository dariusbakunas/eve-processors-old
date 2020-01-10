package db

import (
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

func (d *DB) CleanupJobLogs(category string, characterID int64) {
	rows, err := squirrel.
		Select("id").
		From("jobLogs").
		Where(squirrel.Eq{"category": category}).
		Where(squirrel.Eq{"characterId": characterID}).
		OrderBy("createdAt DESC").Limit(4).
		RunWith(d.db).
		Query()

	if err != nil {
		log.Printf("Failed cleaning up job logs: %v", err)
	}

	var ids []int

	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			log.Printf("Failed cleaning up job logs: %v", err)
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
		log.Printf("Failed cleaning up job logs: %v", err)
	}
}