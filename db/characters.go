package db

import (
	"fmt"
	"github.com/dariusbakunas/eve-processors/utils"
)

import sq "github.com/Masterminds/squirrel"

type Character struct {
	ID           int64
	AccessToken  string
	RefreshToken string
	Expires      int
	Scopes       string
}

func (d *DB) GetCharacters() ([]Character, error) {
	rows, err := sq.Select("id, accessToken, refreshToken, expiresAt, scopes").From("characters").RunWith(d.db).Query()

	if err != nil {
		return nil, fmt.Errorf("sq.Select: %v", err)
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
			return nil, fmt.Errorf("d.crypt.Decrypt: %v", err)
		}

		character.RefreshToken, err = d.crypt.Decrypt(character.RefreshToken)

		if err != nil {
			return nil, fmt.Errorf("d.crypt.Decrypt: %v", err)
		}

		characters = append(characters, character)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err: %v", err)
	}

	return characters, nil
}

func (d *DB) UpdateCharacterTokens(accessToken string, refreshToken string, expiresIn int64, characterId int64) error {
	timestamp := utils.GetCurrentTimestamp()

	encryptedAccessToken, err := d.crypt.Encrypt(accessToken)

	if err != nil {
		return fmt.Errorf("d.crypt.Encrypt: %v", err)
	}

	encryptedRefreshToken, err := d.crypt.Encrypt(refreshToken)

	if err != nil {
		return fmt.Errorf("d.crypt.Encrypt: %v", err)
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