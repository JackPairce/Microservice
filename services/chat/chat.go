package chat

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

// mustEmbedUnimplementedChatServiceServer implements ChatServiceServer.
func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	panic("unimplemented")
}

func (s *Server) SendMessage(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Received message body from client: %s", in.Body)
	return &Message{Body: "Hello From the Server!"}, nil
}
