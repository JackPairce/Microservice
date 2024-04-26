package main

import (
	"log"
	"net"

	"github.com/JackPairce/MicroService/services/fileindexing"
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

	var ListUsers []superpeer.User
	s := superpeer.Server{Users: &ListUsers, Indexer: &fileindexing.FileIndexing{}}
	grpcServer := grpc.NewServer(grpc.StatsHandler(&s))
	superpeer.RegisterSuperPeerServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
