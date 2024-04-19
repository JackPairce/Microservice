package superpeer

import (
	context "context"
	"log"
	"math"
	"math/rand/v2"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
}

func (s *Server) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	// Validate input parameters
	if in.Name == "" || in.Address == "" || in.Password == "" {
		return &RegisterResponse{
				Success: false,
				Id:      int32(math.NaN()),
			},
			status.Errorf(codes.InvalidArgument, "Missing required fields")
	}
	// If the user is already registered, return an error
	// TODO: Implement logic to check if the user is already registered
	if false {
		return &RegisterResponse{
			Success: false,
			Id:      int32(math.NaN()),
		}, status.Errorf(codes.AlreadyExists, "User already registered")
	}
	log.Printf("%s with %s is Registred", in.Name, in.Address)
	return &RegisterResponse{
		Success: true,
		Id:      int32(rand.IntN(999999)),
	}, nil
}

func (s *Server) Login(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	// Validate input parameters
	if in.Address == "" || in.Password == "" {
		return &RegisterResponse{
				Success: false,
				Id:      int32(math.NaN()),
			},
			status.Errorf(codes.InvalidArgument, "Missing required fields")
	}
	// TODO : check if the user is registered or not
	return &RegisterResponse{
		Success: true,
		Id:      int32(rand.IntN(999999)),
	}, nil
}

func (s *Server) SearchFiles(ctx context.Context, in *SearchFilesRequest) (*SearchFilesResponse, error) {
	// Validate input parameters
	if in.Query == "" {
		return &SearchFilesResponse{},
			nil
	}
	return &SearchFilesResponse{}, nil
}

func (s *Server) GetPeerFiles(ctx context.Context, in *FileList) (*Empty, error) {
	// Validate input parameters
	if len(in.Files) == 0 {
		return &Empty{}, status.Errorf(codes.InvalidArgument, "Missing required field: peerAddress")
	}
	return &Empty{}, nil
}

func (s *Server) mustEmbedUnimplementedSuperPeerServer() {}
