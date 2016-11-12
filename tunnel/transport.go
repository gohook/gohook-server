package tunnel

import (
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/gohook/gohook-server/pb"
	"github.com/gohook/gohook-server/user"
	"github.com/satori/go.uuid"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"time"
)

type GohookTunnelServer struct {
	// Auth Service
	auth user.AuthService

	// Queue for getting notified when hooks come in
	queue HookQueue

	// Session Store for adding new sessions
	sessions *SessionStore

	// Message logger
	logger log.Logger
}

func (s GohookTunnelServer) SendToStream(accountId user.AccountId, message HookCall) error {
	sessions, err := s.sessions.FindByAccountId(accountId)
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

	token, err := getTokenFromContext(streamCtx)
	if err != nil {
		return err
	}

	account, err := s.auth.AuthAccountFromToken(token)
	if err != nil {
		return err
	}

	s.logger.Log("msg", "Have authed user", "account_id", account.Id, "account_token", account.Token)

	id := uuid.NewV4()
	newSession := &Session{
		Id:        SessionId(id.String()),
		AccountId: account.Id,
		Start:     time.Now(),
		Stream:    stream,
	}

	err = s.sessions.Add(newSession)
	if err != nil {
		return err
	}
	s.logger.Log("msg", "Added stream to list", "streamId", newSession.Id, "account_id", newSession.AccountId)

	for {
		select {
		case <-streamCtx.Done():
			err := streamCtx.Err()
			s.logger.Log("msg", "Stream done", "sessionId", newSession.Id, "err", err)
			return s.sessions.Remove(newSession.AccountId, newSession.Id)
		}

	}
}

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", errors.New("Missing context data from stream")
	}

	mdToken, ok := md["token"]
	if !ok && len(mdToken) > 0 {
		return "", errors.New("Missing auth token in GRPC request")
	}

	return mdToken[0], nil
}

func MakeTunnelServer(authService user.AuthService, q HookQueue, logger log.Logger) (*GohookTunnelServer, error) {
	queuec, err := q.Listen()
	if err != nil {
		return nil, err
	}

	sessions := NewSessionStore()

	server := &GohookTunnelServer{
		auth:     authService,
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
					return
				}

				logger.Log("msg", "Handling incoming messsage...", "message", msg.Hook.Id)
				server.SendToStream(msg.AccountId, msg.Hook)
			}

		}
	}()

	return server, nil
}
