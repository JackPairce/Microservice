package superpeer

import (
	"context"
	"log"
	"math"

	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
)

func (s *Server) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

// HandleRPC processes the RPC stats.
func (s *Server) HandleRPC(ctx context.Context, sts stats.RPCStats) {
	// log.Println("HandleRPC")
}

func (s *Server) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	log.Println("Peer connected:", info.RemoteAddr.String())
	// ctx = context.WithValue(ctx, info.RemoteAddr.String(), "nil")
	s.mu.Lock()
	(*s.ActivePeers)[info.RemoteAddr.String()] = int64(math.NaN())
	s.mu.Unlock()
	return ctx
}

// HandleConn processes the Conn stats.
func (s *Server) HandleConn(ctx context.Context, sts stats.ConnStats) {
	switch sts.(type) {
	case *stats.ConnEnd:
		p, _ := peer.FromContext(ctx)
		// peerId := ctx.Value(p.Addr.String())
		s.mu.Lock()
		peerId := (*s.ActivePeers)[p.Addr.String()]
		delete(*s.ActivePeers, p.Addr.String())
		log.Println("Peer disconnected: ", peerId)
		s.mu.Unlock()
		if err := s.DB.UserLogout(peerId); err != nil {
			log.Println("Error while logging out user: ", err)
		}
	}
}
