package tunnel

import (
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
	// NOTE: SessionID must be unique per instance, but trackable with the given user. Maybe
	// set up a table with userId -> sessionId mapping
	sessions map[SessionID]pb.Gohook_TunnelServer

	// Message logger
	logger log.Logger
}

func (s GohookTunnelServer) SendToStream(id SessionID, message *pb.HookCall) error {
	if stream, ok := s.sessions[id]; ok {
		err := stream.Send(&pb.TunnelResponse{
			Event: &pb.TunnelResponse_Hook{
				Hook: message,
			},
		})
		return err
	}

	// It's not an error if we don't have the stream (it probably lives on a different process)
	return nil
}

// Tunnel transport handler
func (s *GohookTunnelServer) Tunnel(req *pb.TunnelRequest, stream pb.Gohook_TunnelServer) error {
	// TODO: Need to generate a unique ID for this session and link it to the user's ID

	tickChan := time.NewTicker(time.Second * 5).C
	streamCtx := stream.Context()

	s.logger.Log("msg", "Added stream to list", "streamId", req.Id)
	s.sessions[SessionID(req.Id)] = stream

	for {
		select {
		case <-streamCtx.Done():
			err := streamCtx.Err()
			s.logger.Log("msg", "Stream done", "err", err)
			delete(s.sessions, SessionID(req.Id))
			return nil
		case <-tickChan:
			// Case for testing. Use broadcast to trigger a message
			err := s.queue.Broadcast("Hello")
			if err != nil {
				s.logger.Log("msg", "Failed to broadcast", "error", err)
			}
		}

	}
}

func MakeTunnelServer(q gohookd.HookQueue, logger log.Logger) (*GohookTunnelServer, error) {
	queuec, err := q.Listen()
	if err != nil {
		return nil, err
	}

	server := &GohookTunnelServer{
		logger:   logger,
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
					continue
				}

				logger.Log("msg", "Handling incoming messsage...", "message", msg)
				server.SendToStream("myid", &pb.HookCall{
					Id: msg.(string),
				})
			}

		}
	}()

	return server, nil
}
