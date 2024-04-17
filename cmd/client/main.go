package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/JackPairce/MicroService/services/chat"
)

func main() {
	port := "8080"
	address := "localhost"
	// Set up a connection to the server.
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := chat.NewChatServiceClient(conn)

	// Contact the server and print out its response.
	message := chat.Message{Body: "Hello From the Client!"}
	response, err := c.SendMessage(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling SendMessage: %s", err)
	}
	log.Printf("Response from the server: %s", response.Body)
}
