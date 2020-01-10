package db

import (
	"database/sql"
	"fmt"
	"github.com/dariusbakunas/eve-processors/utils"
	"os"
)

type DB struct {
	db *sql.DB
	crypt *utils.Crypt
}

func (d *DB) Close() error {
	return d.db.Close()
}

func NewDB(connection string, database string, username string, password string, tokenSecret string) (*DB, error) {
	var err error
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, password, connection, database))

	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	crypt := &utils.Crypt{Key: tokenSecret}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

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

func (d *DB) Encrypt(plainText string) (string, error) {
	return d.crypt.Encrypt(plainText)
}

func (d *DB) Decrypt(cipherText string) (string, error) {
	return d.crypt.Decrypt(cipherText);
}