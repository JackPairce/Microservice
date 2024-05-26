package superpeer

import (
	"context"
	"math"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
)

func (s *Server) UnaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	if !(info.FullMethod == "/protos.SuperPeer/Register" || info.FullMethod == "/protos.SuperPeer/Login" || info.FullMethod == "/protos.SuperPeer/Ping") {
		p, _ := peer.FromContext(ctx)
		s.mu.Lock()
		defer s.mu.Unlock()
		if (*s.ActivePeers)[p.Addr.String()] == int64(math.NaN()) {
			return nil, status.Errorf(codes.Internal, "Peer not registered")
		}
	}

	return handler(ctx, req)
}
