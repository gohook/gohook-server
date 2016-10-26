package tunnel

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
type QueueMessage struct {
	SessionId SessionID
}

type ReceiveC chan *QueueMessage

type HookQueue interface {
	Broadcast(message *QueueMessage) error
	Listen() (ReceiveC, error)
}
