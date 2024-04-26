package peer

import (
	"io"
	"os"
	"path"

	types "github.com/JackPairce/MicroService/services/types"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Peer struct {
	Path string
	Name string
}

const chunksize = 1024 * 3

func (p *Peer) SendFile(in *types.File, stream Peer_SendFileServer) error {
	// "/home/jackpairce/Documents/"+in.Name
	f, err := os.OpenFile(path.Join(p.Path, p.Name, in.Name), os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return status.Error(codes.Internal, "file not exist")
		}
		return status.Errorf(codes.Internal, "Failed to open file: %v", err)
	}
	defer f.Close()
	// check if is not a directory
	if stat, err := f.Stat(); err != nil {
		return status.Errorf(codes.Internal, "Failed to stat file: %v", err)
	} else if stat.IsDir() {
		return status.Errorf(codes.Internal, "Failed to send file: %v", "is a directory")
	}
	chunk := &types.FileData{
		Data: make([]byte, chunksize),
	}
	var n int
Loop:
	for {
		n, err = f.Read(chunk.Data)
		switch err {
		case nil:
		case io.EOF:
			break Loop
		default:
			return status.Errorf(codes.Internal, "Failed to read file: %v", err)
		}
		chunk.Data = chunk.Data[:n]
		if err := stream.Send(chunk); err != nil {
			return err
		}
	}
	return nil
}
func (p *Peer) mustEmbedUnimplementedPeerServer() {
	panic("implement me")
}
