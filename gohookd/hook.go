package gohookd

type HookID string

type HookList []*Hook

type Hook struct {
	Id     HookID `json:"id"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

type HookCall struct {
	Id     HookID `json:"id"`
	Method string `json:"method"`
	Body   string `json:"body"`
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

/*
Queue Interface
---------------

The Queue interface describes a system for broadcasting
and receiving messages from any other running gohookd
process. It acts as a message broker for allowing all
processes to know about incoming hook messages and
allows the process with the connected client to handle
sending the message down to the client.
*/
type QueueMessage interface{}

type ReceiveC chan QueueMessage

type HookQueue interface {
	Broadcast(message QueueMessage) error
	Listen() (ReceiveC, error)
}
