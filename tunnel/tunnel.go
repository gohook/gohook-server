package tunnel

import (
	"fmt"
	"github.com/gohook/gohook-server/pb"
	"time"
)

type GohookTunnelServer struct {
	// Keeps a list of all the sessions and the session ids for sending to correct client
	sessions map[string]pb.Gohook_TunnelServer
}

func MakeTunnelServer() *GohookTunnelServer {
	return &GohookTunnelServer{}
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
