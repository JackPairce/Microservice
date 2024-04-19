package chat

import (
	"bufio"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
)

type Server struct {
}

// mustEmbedUnimplementedChatServiceServer implements ChatServiceServer.
func (s *Server) mustEmbedUnimplementedChatServiceServer() {
	panic("unimplemented")
}

func (s *Server) SendMessage(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Client: %s", in.Body)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return &Message{Body: text}, nil
}
