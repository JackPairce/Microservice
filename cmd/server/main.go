package main

import (
	"log"
	"net"

	"github.com/JackPairce/MicroService/services/database"
	"github.com/JackPairce/MicroService/services/superpeer"
	"google.golang.org/grpc"
)

func main() {
	port := "8080"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on port %s", port)

	db := database.DB{}
	err = db.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	ActivePeers := make(map[string]int64)

	s := superpeer.Server{
		DB:          &db,
		ActivePeers: &ActivePeers,
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(&s),
		grpc.UnaryInterceptor(s.UnaryInterceptor),
	)
	superpeer.RegisterSuperPeerServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
