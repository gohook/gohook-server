package tunnel

import (
	"github.com/gohook/gohook-server/user"
)

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

type HookCall struct {
	Id     string `json:"id"`
	Method string `json:"method"`
	Body   []byte `json:"body"`
}

type QueueMessage struct {
	AccountId user.AccountId
	Hook      HookCall
}

type ReceiveC chan *QueueMessage

type HookQueue interface {
	Broadcast(message *QueueMessage) error
	Listen() (ReceiveC, error)
}
