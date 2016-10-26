package tunnel

import (
	"github.com/go-kit/kit/log"
	"github.com/gohook/gohook-server/gohookd"
	"github.com/gohook/gohook-server/pb"
	"github.com/satori/go.uuid"
	"time"
)

type GohookTunnelServer struct {
	// Queue for getting notified when hooks come in
	queue HookQueue

	// Session Store for adding new sessions
	sessionStore SessionStore

	// Keeps a list of all the sessions and the session ids for sending to correct client
	// NOTE: SessionID must be unique per instance, but trackable with the given user. Maybe
	// set up a table with userId -> sessionId mapping
	sessions map[SessionID]pb.Gohook_TunnelServer

	// Message logger
	logger log.Logger
}

func (s GohookTunnelServer) SendToStream(id SessionID, message gohookd.HookCall) error {
	if stream, ok := s.sessions[id]; ok {
		err := stream.Send(&pb.TunnelResponse{
			Event: &pb.TunnelResponse_Hook{
				Hook: &pb.HookCall{
					Id:   string(message.Id),
					Body: message.Body,
				},
			},
		})
		return err
	}

	// It's not an error if we don't have the stream (it probably lives on a different process)
	return nil
}

// Tunnel transport handler
func (s *GohookTunnelServer) Tunnel(req *pb.TunnelRequest, stream pb.Gohook_TunnelServer) error {
	streamCtx := stream.Context()
	id := uuid.NewV4()
	newSession := &Session{
		Id:     SessionID(id.String()),
		UserId: req.Id,
		Start:  time.Now(),
	}

	err := s.sessionStore.Add(newSession)
	if err != nil {
		return err
	}

	s.logger.Log("msg", "Added stream to list", "streamId", newSession.Id, "userId", newSession.UserId)
	s.sessions[newSession.Id] = stream

	for {
		select {
		case <-streamCtx.Done():
			err := streamCtx.Err()
			s.logger.Log("msg", "Stream done", "sessionId", newSession.Id, "err", err)
			delete(s.sessions, newSession.Id)
			return s.sessionStore.Remove(newSession.Id)
		}

	}
}

func MakeTunnelServer(sessions SessionStore, q HookQueue, logger log.Logger) (*GohookTunnelServer, error) {
	queuec, err := q.Listen()
	if err != nil {
		return nil, err
	}

	server := &GohookTunnelServer{
		logger:       logger,
		sessionStore: sessions,
		queue:        q,
		sessions:     make(map[SessionID]pb.Gohook_TunnelServer),
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
				server.SendToStream(msg.SessionId, msg.Hook)
			}

		}
	}()

	return server, nil
}
