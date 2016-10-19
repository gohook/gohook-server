package gohookd

type HookStore interface {
	Add(hook *Hook) error
	Remove(hookId HookID) (*Hook, error)
	Find(hookId HookID) (*Hook, error)
	FindAll() (HookList, error)
}
