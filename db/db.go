package db

import (
	"database/sql"
	"fmt"
	"github.com/dariusbakunas/eve-processors/utils"
	"os"
)

import sq "github.com/Masterminds/squirrel"

type DB struct {
	db *sql.DB
	crypt *utils.Crypt
}

func (d *DB) Close() error {
	return d.db.Close()
}

func NewDB(connection string, database string, username string, password string, tokenSecret string) (*DB, error) {
	var err error
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", username, password, connection, database))

	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	crypt := &utils.Crypt{Key: tokenSecret}

	return &DB {
		db: db,
		crypt: crypt,
	}, nil
}

func InitializeDb() (*DB, error) {
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		return nil, fmt.Errorf("TOKEN_SECRET must be set")
	}

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

	db, err := NewDB(dbConnection, dbDatabase, dbUsername, dbPassword, tokenSecret)

	if err != nil {
		return nil, fmt.Errorf("NewDB: %v", err)
	}

	return db, nil
}


type Character struct {
	ID           int64
	AccessToken  string
	RefreshToken string
	Expires      int
	Scopes       string
}

func (d *DB) Encrypt(plainText string) (string, error) {
	return d.crypt.Encrypt(plainText)
}

func (d *DB) Decrypt(cipherText string) (string, error) {
	return d.crypt.Decrypt(cipherText);
}

func (d *DB) GetCharacters() ([]Character, error) {
	rows, err := sq.Select("id, accessToken, refreshToken, expiresAt, scopes").From("characters").RunWith(d.db).Query()

	if err != nil {
		return nil, fmt.Errorf("mysql: %v", err)
	}

	defer rows.Close()

	var characters []Character

	for rows.Next() {
		var character Character

		err := rows.Scan(&character.ID, &character.AccessToken, &character.RefreshToken, &character.Expires, &character.Scopes)

		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %v", err)
		}

		character.AccessToken, err = d.crypt.Decrypt(character.AccessToken)

		if err != nil {
			return nil, fmt.Errorf("db.crypt.Decrypt: %v", err)
		}

		character.RefreshToken, err = d.crypt.Decrypt(character.RefreshToken)

		if err != nil {
			return nil, fmt.Errorf("db.crypt.Decrypt: %v", err)
		}

		characters = append(characters, character)
	}

	return characters, nil
}

func (d *DB) UpdateCharacterTokens(accessToken string, refreshToken string, expiresIn int64, characterId int64) error {
	timestamp := utils.GetCurrentTimestamp()

	encryptedAccessToken, err := d.crypt.Encrypt(accessToken)

	if err != nil {
		return fmt.Errorf("db.crypt.Encrypt: %v", err)
	}

	encryptedRefreshToken, err := d.crypt.Encrypt(refreshToken)

	if err != nil {
		return fmt.Errorf("db.crypt.Encrypt: %v", err)
	}

	_, err = sq.Update("characters").
		Set("accessToken", encryptedAccessToken).
		Set("refreshToken", encryptedRefreshToken).
		Set("expiresAt", expiresIn * 1000 + timestamp).
		Where(sq.Eq{"id": characterId}).
		RunWith(d.db).
		Exec()

	if err != nil {
		return fmt.Errorf("sq.Update: %v", err)
	}

	return nil
}