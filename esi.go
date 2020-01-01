package esi

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)

import _ "github.com/go-sql-driver/mysql"
import sq "github.com/Masterminds/squirrel"

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func initializeDb() (*sql.DB, error) {
	dbConnection := os.Getenv("DB_CONNECTION")

	if dbConnection == "" {
		return nil, fmt.Errorf("DB_CONNECTION must be set")
	}

	dbDatabase := os.Getenv("DB_DATABASE")

	if dbDatabase == "" {
		return nil, fmt.Errorf("DB_DATABASE must be set")
	}

	dbUsername := os.Getenv("DB_USERNAME")

	if dbUsername == "" {
		return nil, fmt.Errorf("DB_USERNAME must be set")
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD must be set")
	}

	var err error
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", dbUsername, dbPassword, dbConnection, dbDatabase))

	if err != nil {
		return nil, fmt.Errorf("mysql: %v", err)
	}

	return db, nil
}

type Character struct {
	id  int
	accessToken string
	refreshToken string
	expires int
}

func getCharacters(db *sql.DB) ([]Character, error) {
	rows, err := sq.Select("id, accessToken, refreshToken, expiresAt").From("characters").RunWith(db).Query()

	if err != nil {
		return nil, fmt.Errorf("mysql: %v", err)
	}

	defer rows.Close()

	var characters []Character

	for rows.Next() {
		var character Character

		err := rows.Scan(&character.id, &character.accessToken, &character.refreshToken, &character.expires)

		if err != nil {
			return nil, fmt.Errorf("mysql: %v", err)
		}

		characters = append(characters, character)
	}

	return characters, nil
}

func Process() {
	db, err := initializeDb()

	if err != nil {
		log.Fatalf("initializeDb: %v", err)
	}

	characters, err := getCharacters(db)

	if err != nil {
		log.Fatal(err)
	}

	for _, character := range characters {
		log.Printf("ID: %d", character.id)
	}

	defer db.Close()
}

func Esi(ctx context.Context, m PubSubMessage) error {
	Process()

	return nil
}