package superpeer

import (
	context "context"
	"log"
	"math"
	sync "sync"

	"github.com/JackPairce/MicroService/services/database"
	t "github.com/JackPairce/MicroService/services/types"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
)

type Server struct {
	DB *database.DB
	UnimplementedSuperPeerServer

	mu          sync.Mutex
	ActivePeers *map[string]int64
}

func (s *Server) UserConnection(ctx context.Context, in *RegisterRequest, IsRegister bool) (*RegisterResponse, error) {
	// Validate the fields
	if in.Peeeraddress == "" || in.Name == "" || in.Password == "" {
		return &RegisterResponse{
				Success: false,
				Id:      int64(math.NaN()),
			},
			status.Errorf(codes.InvalidArgument, "Missing required fields")
	}

	// If the user is already registered, return an error
	err := s.DB.CheckUserExistence(&t.User{Name: in.Name})
	if err != nil {
		return &RegisterResponse{
				Success: false,
				Id:      int64(math.NaN()),
			},
			status.Errorf(codes.Internal, err.Error())
	}

	var ID int64
	if IsRegister {
		// Create a new user
		ID, err = s.DB.AddUser(&t.User{
			Name:        in.Name,
			Password:    in.Password,
			Peeraddress: in.Peeeraddress,
		})
		if err != nil {
			return &RegisterResponse{
					Success: false,
					Id:      int64(math.NaN()),
				},
				status.Errorf(codes.Internal, err.Error())
		}
		log.Printf("%s with %s is Registred", in.Name, in.Peeeraddress)
	} else {
		// check if the password is correct
		ID, err = s.DB.CheckUserPassword(&t.User{Name: in.Name, Password: in.Password})
		if err != nil {
			return &RegisterResponse{
					Success: false,
					Id:      int64(math.NaN()),
				},
				status.Errorf(codes.Internal, err.Error())
		}
	}
	// save Peer address in the database
	err = s.DB.UserLogin(ID, in.Peeeraddress)
	if err != nil {
		return &RegisterResponse{
				Success: false,
				Id:      int64(math.NaN()),
			},
			status.Errorf(codes.Internal, err.Error())
	}

	// save peer id in the ActivePeers map
	p, _ := peer.FromContext(ctx)
	s.mu.Lock()
	(*s.ActivePeers)[p.Addr.String()] = ID
	s.mu.Unlock()

	// return Success
	return &RegisterResponse{
		Success: true,
		Id:      ID,
	}, nil
}

func (s *Server) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	return s.UserConnection(ctx, in, true)
}

func (s *Server) Login(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	return s.UserConnection(ctx, in, false)
}

func (s *Server) SearchFiles(ctx context.Context, in *SearchFilesRequest) (*SearchFilesResponse, error) {
	files, err := s.DB.SearchFile(int(in.Id), in.Filename, false)
	if err != nil {
		return &SearchFilesResponse{
			Results: &t.FileList{Files: []*t.File{}},
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &SearchFilesResponse{
		Results: &t.FileList{Files: *files},
	}, nil
}

func (s *Server) GetPeerFiles(ctx context.Context, in *t.FileList) (*Empty, error) {
	// Validate input parameters
	if len(in.Files) == 0 {
		return &Empty{}, status.Errorf(codes.InvalidArgument, "There is no files")
	}
	// Add files to the database
	for _, file := range in.Files {
		s.DB.AddFile(file)
	}

	if AllFiles, err := s.DB.GetAllFiles(int(in.Files[0].Ownerid)); err != nil {
		return &Empty{}, status.Errorf(codes.Internal, err.Error())
	} else if len(in.Files) != len(*AllFiles) {
		for _, f := range *AllFiles {
			if !func() bool {
				for _, file := range in.Files {
					if file.Id == f.Id {
						return true
					}
				}
				return false
			}() {
				if err := s.DB.DeleteFile(int(f.Id)); err != nil {
					return &Empty{}, status.Errorf(codes.Internal, err.Error())
				}
			}
		}
	}
	return &Empty{}, nil
}

func (s *Server) GetPeerConnexion(ctx context.Context, in *PeerId) (*PeerConnexion, error) {
	peeraddress, err := s.DB.GetUserPeerAdress(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &PeerConnexion{
		Peeraddress: peeraddress,
	}, nil
}
