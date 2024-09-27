package entities

import "time"

type AccessToken struct {
	Token string
	IP    string
	Exp   time.Time
}

type RefreshToken struct {
	TokenHash string
	UserID    int
	IP        string
	Exp       time.Time
}
