package superpeer

import (
	context "context"
	"log"
	"math"
	"math/rand/v2"
	"slices"

	f "github.com/JackPairce/MicroService/services/fileindexing"
	t "github.com/JackPairce/MicroService/services/types"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	Users   *[]User
	Indexer *f.FileIndexing
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
	ID := int32(rand.IntN(999999))
	*s.Users = append(*s.Users, User{Id: ID})
	println("reg", s.Users)
	return &RegisterResponse{
		Success: true,
		Id:      ID,
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
	if in.Filename == "" {
		return &SearchFilesResponse{
			Results: &t.FileList{Files: s.Indexer.GetUniqueFileNames(in.Id)},
		}, nil
	}
	// filter the search results to not include files that are already owned by the requesting user
	var filteredFiles []*t.File
	for _, file := range s.Indexer.SearchFiles(in.Filename) {
		if !slices.Contains(file.Ownerid, in.Id) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return &SearchFilesResponse{
		Results: &t.FileList{Files: filteredFiles},
	}, nil
}

func (s *Server) GetPeerFiles(ctx context.Context, in *t.FileList) (*Empty, error) {
	// Validate input parameters
	if len(in.Files) == 0 {
		return &Empty{}, status.Errorf(codes.InvalidArgument, "There is no files")
	}
	log.Println("Files", in.Files)

	// adding missing files to the server's file list
	for _, file := range in.Files {
		FILES := s.Indexer.SearchFiles(file.Name)
		if len(FILES) == 0 {
			s.Indexer.AddFile(file)
		} else if len(FILES) == 1 && !slices.Contains(FILES[0].Ownerid, file.Ownerid[0]) {
			s.Indexer.UpdateFile(file)
		} else {
			log.Fatalln("Too many files with the same name found. This should not happen.")
		}
	}
	return &Empty{}, nil
}

func (s *Server) mustEmbedUnimplementedSuperPeerServer() {}
