package tunnel

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/pb"
	"time"
)

type SessionID string

type GohookTunnelServer struct {
	// Queue for getting notified when hooks come in
	queue gohookd.HookQueue

	// Keeps a list of all the sessions and the session ids for sending to correct client
	sessions map[SessionID]pb.Gohook_TunnelServer
}

func MakeTunnelServer(q gohookd.HookQueue, logger log.Logger) (*GohookTunnelServer, error) {
	queuec, err := q.Listen()
	if err != nil {
		return nil, err
	}

	server := &GohookTunnelServer{
		queue:    q,
		sessions: make(map[SessionID]pb.Gohook_TunnelServer),
	}

	// Process for handling queue messages
	go func() {
		for {
			select {

			// process incoming messages from RedisPubSub, and send messages.
			case msg := <-queuec:
				if msg == nil {
					logger.Log("msg", "Message Channel has closed. Exiting.")
				}

				logger.Log("msg", "Handling incoming messsage...", "message", msg)
			}

		}
	}()

	return server, nil
}

// Tunnel transport handler
func (s *GohookTunnelServer) Tunnel(req *pb.TunnelRequest, stream pb.Gohook_TunnelServer) error {
	// pass the stream and request data down to the service?
	// that way we can handle the logic there... but... then we are mixing transport and
	// business logic...
	// is there some way we can expose an interface to the service?
	for {
		err := stream.Send(&pb.TunnelResponse{
			Event: &pb.TunnelResponse_Hook{
				Hook: &pb.HookCall{
					Id: "Hello, World",
				},
			},
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
}
