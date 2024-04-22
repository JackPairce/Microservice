package superpeer

import (
	"context"
	"log"

	"google.golang.org/grpc/stats"
)

func (s *Server) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	// log.Print("TagRPC")
	return context.Background()
}

// HandleRPC processes the RPC stats.
func (s *Server) HandleRPC(context.Context, stats.RPCStats) {
	// log.Println("HandleRPC")
}

func (s *Server) TagConn(context.Context, *stats.ConnTagInfo) context.Context {

	// log.Println("Tag Conn")
	return context.Background()
}

// HandleConn processes the Conn stats.
func (s *Server) HandleConn(c context.Context, sts stats.ConnStats) {
	switch sts.(type) {
	case *stats.ConnBegin:
		log.Println("ConnBegin")
		if len(*s.Users) != 0 {
			println("Available Users")
			for _, user := range *s.Users {
				println("-", user.Id)
			}
		}
	case *stats.ConnEnd:
		log.Println("ConnEnd ")
		//fmt.Printf("client %d disconnected", s.userIdMap[ctx.Value("user_counter")])
		return
	}
}
