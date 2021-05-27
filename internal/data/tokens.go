package data

import "database/sql"

type AuthToken struct {
	Service      string
	AccessToken  string
	RefreshToken string
}

type AuthTokenModel struct {
	DB *sql.DB
}

func (auth *AuthTokenModel) Update(t AuthToken) error {
	stm := `UPDATE auth_tokens SET access_token = $1, refresh_token = $2
			WHERE service = $3`
	args := []interface{}{t.AccessToken, t.RefreshToken, t.Service}
	_, err := auth.DB.Exec(stm, args...)
	return err
}

func (auth *AuthTokenModel) Get(srv string) (*AuthToken, error) {
	t := &AuthToken{}
	stm := `SELECT access_token, refresh_token FROM tokens
			WHERE service = $1`

	row := auth.DB.QueryRow(stm, srv)
	err := row.Scan(&t.AccessToken, &t.RefreshToken)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (auth *AuthTokenModel) Insert(srv string) error {
	stm := `INSERT INTO TOKENS (service, access_token, refresh_token)
			values ($1, '', '')`
	_, err := auth.DB.Exec(stm, srv)
	return err
}
