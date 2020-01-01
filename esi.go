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

func Esi(ctx context.Context, m PubSubMessage) error {
	connectionName := os.Getenv("CLOUD_SQL_CONNECTION_NAME")

	if connectionName == "" {
		return fmt.Errorf("CLOUD_SQL_CONNECTION_NAME must be set")
	}

	dbUsername := os.Getenv("DB_USERNAME")

	if dbUsername == "" {
		return fmt.Errorf("DB_USERNAME must be set")
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		return fmt.Errorf("DB_PASSWORD must be set")
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/eve_gql_pd", dbUsername, dbPassword, connectionName))

	if err != nil {
		return fmt.Errorf("mysql: %v", err)
	}

	rows, err := sq.Select("id, name").From("characters").RunWith(db).Query()

	if err != nil {
		return fmt.Errorf("mysql: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id int
			name string
		)

		err := rows.Scan(&id, &name)

		if err != nil {
			return fmt.Errorf("mysql: %v", err)
		}

		log.Printf("Name: %s", name);
	}

	defer db.Close()

	return nil
}