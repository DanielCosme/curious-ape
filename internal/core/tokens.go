package core

type Token struct {
	Service      string `json:"-"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthTokenModel interface {
	Update(t Token) error
	Get(service string) (*Token, error)
	Insert(service string) error
}
