package user

type AccountId string

type AccountToken string

type Account struct {
	Id    AccountId
	Token AccountToken
}

type AccountStore interface {
	Add(account *Account) error
	Remove(accountId AccountId) (*Account, error)
	Find(accountId AccountId) (*Account, error)
	FindByToken(token AccountToken) (*Account, error)
}
