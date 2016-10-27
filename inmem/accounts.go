package inmem

import (
	"errors"
	"sync"

	"github.com/gohook/gohook-server/user"
)

type InMemAccounts struct {
	mtx      sync.RWMutex
	accounts map[user.AccountId]*user.Account
}

func NewInMemAccounts() user.AccountStore {
	accounts := make(map[user.AccountId]*user.Account)
	accounts["myid"] = &user.Account{
		Id:    user.AccountId("myid"),
		Token: user.AccountToken("brianegizi"),
	}
	return &InMemAccounts{
		accounts: accounts,
	}
}

func (i *InMemAccounts) Add(a *user.Account) error {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	i.accounts[a.Id] = a
	return nil
}

func (i *InMemAccounts) Remove(id user.AccountId) (*user.Account, error) {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	account, ok := i.accounts[id]
	if ok {
		delete(i.accounts, id)
		return account, nil
	}
	return nil, errors.New("Not Found")
}

func (i *InMemAccounts) Find(id user.AccountId) (*user.Account, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	if val, ok := i.accounts[id]; ok {
		return val, nil
	}
	return nil, errors.New("Not Found")
}

func (i *InMemAccounts) FindByToken(token user.AccountToken) (*user.Account, error) {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	for _, account := range i.accounts {
		if account.Token == token {
			return account, nil
		}
	}
	return nil, errors.New("Not Found")
}
