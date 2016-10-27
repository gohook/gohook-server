package gohookd

type HookID string

type HookList []*Hook

type Hook struct {
	Id        HookID `json:"id"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	AccountId string `json:"account_id"`
}

type HookRequest struct {
	Method string `json:"method"`
}

// HookStore is an interface defining the methods used to store hooks
type HookStore interface {
	Add(hook *Hook) error
	Remove(hookId HookID) (*Hook, error)
	Find(hookId HookID) (*Hook, error)
	FindAll() (HookList, error)
}
