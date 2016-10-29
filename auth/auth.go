package auth

import (
	"github.com/gohook/gohook-server/user"
)

type basicAuthService struct {
	accounts user.AccountStore
}

func NewAuthService(accounts user.AccountStore) user.AuthService {
	return &basicAuthService{
		accounts: accounts,
	}
}

func (a basicAuthService) AuthAccountFromToken(token string) (*user.Account, error) {
	return a.accounts.FindByToken(user.AccountToken(token))
}
