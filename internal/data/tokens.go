package data

import (
	"database/sql"

	"github.com/danielcosme/curious-ape/internal/auth"
)

type AuthTokenModel struct {
	DB *sql.DB
}

func (a *AuthTokenModel) Update(t auth.Token) error {
	stm := `UPDATE auth_tokens SET access_token = $1, refresh_token = $2
			WHERE service = $3`
	args := []interface{}{t.AccessToken, t.RefreshToken, t.Service}
	_, err := a.DB.Exec(stm, args...)
	return err
}

func (a *AuthTokenModel) Get(srv string) (*auth.Token, error) {
	t := &auth.Token{}
	stm := `SELECT access_token, refresh_token FROM tokens
			WHERE service = $1`

	row := a.DB.QueryRow(stm, srv)
	err := row.Scan(&t.AccessToken, &t.RefreshToken)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (a *AuthTokenModel) Insert(srv string) error {
	stm := `INSERT INTO TOKENS (service, access_token, refresh_token)
			values ($1, '', '')`
	_, err := a.DB.Exec(stm, srv)
	return err
}
