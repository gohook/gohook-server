package tunnel

import (
	"github.com/go-kit/kit/log"
	"github.com/gohook/gohook-server/pb"
	"github.com/satori/go.uuid"
	"time"
)

type GohookTunnelServer struct {
	// Queue for getting notified when hooks come in
	queue HookQueue

	// Session Store for adding new sessions
	sessions *SessionStore

	// Message logger
	logger log.Logger
}

func (s GohookTunnelServer) SendToStream(userId string, message HookCall) error {
	sessions, err := s.sessions.FindByUserId(userId)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		session.Stream.Send(&pb.TunnelResponse{
			Event: &pb.TunnelResponse_Hook{
				Hook: &pb.HookCall{
					Id:   string(message.Id),
					Body: message.Body,
				},
			},
		})
	}

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
		Stream: stream,
	}

	err := s.sessions.Add(newSession)
	if err != nil {
		return err
	}
	s.logger.Log("msg", "Added stream to list", "streamId", newSession.Id, "userId", newSession.UserId)

	for {
		select {
		case <-streamCtx.Done():
			err := streamCtx.Err()
			s.logger.Log("msg", "Stream done", "sessionId", newSession.Id, "err", err)
			return s.sessions.Remove(newSession.UserId, newSession.Id)
		}

	}
}

func MakeTunnelServer(q HookQueue, logger log.Logger) (*GohookTunnelServer, error) {
	queuec, err := q.Listen()
	if err != nil {
		return nil, err
	}

	sessions := NewSessionStore()

	server := &GohookTunnelServer{
		logger:   logger,
		queue:    q,
		sessions: sessions,
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
				server.SendToStream(msg.UserId, msg.Hook)
			}

		}
	}()

	return server, nil
}
