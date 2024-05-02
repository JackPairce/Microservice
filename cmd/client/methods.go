package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	p "github.com/JackPairce/MicroService/services/peer"
	s "github.com/JackPairce/MicroService/services/superpeer"
	t "github.com/JackPairce/MicroService/services/types"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
)

type NodeInfo struct {
	Name          string
	Pass          string
	Id            int64
	ctx           s.SuperPeerClient
	SearchedFiles []*t.File
	MyServerPort  string
	localpath     string
	grpcServer    *grpc.Server
}

func (nd *NodeInfo) Register() bool {
	// pass data in the context
	// ctx = metadata.AppendToOutgoingContext(ctx, "id", "123")
	// md, ok := metadata.FromOutgoingContext(ctx)
	// fmt.Printf("After update. Metadata: %v, Found: %t\n", md, ok)
	res, err :=
		nd.
			ctx.
			Register(context.Background(), &s.RegisterRequest{
				Name:         nd.Name,
				Password:     nd.Pass,
				Peeeraddress: GetLocalIP() + ":" + nd.MyServerPort,
			})
	if err != nil {
		er := strings.Split(err.Error(), " = ")
		fmt.Println("Error Register:", er[len(er)-1])
		return false
	}
	fmt.Println("Registration successful! with id: ", res.Id)
	nd.Id = res.Id
	MyID = res.Id
	return res.Success
}
func (nd *NodeInfo) Login() bool {
	res, err := nd.ctx.Login(context.Background(), &s.RegisterRequest{
		Name:         nd.Name,
		Password:     nd.Pass,
		Peeeraddress: GetLocalIP() + ":" + nd.MyServerPort,
	})
	if err != nil {
		er := strings.Split(err.Error(), " = ")
		fmt.Println("Error logging in:", er[len(er)-1])
		return false
	}
	fmt.Println("Login successful!", res.Id)
	nd.Id = res.Id
	MyID = res.Id
	return res.Success
}

func (nd *NodeInfo) SearchFiles(searchTerm string) bool {
	res, err := nd.ctx.SearchFiles(context.Background(), &s.SearchFilesRequest{
		Id:       nd.Id,
		Filename: searchTerm,
	})

	if err != nil {
		fmt.Println("Error searching files:", err)
		return strings.Contains(err.Error(), "Peer not registered")
	}

	if len(res.Results.Files) == 0 {
		fmt.Println("No files found matching the search term.")
	} else {
		nd.SearchedFiles = res.Results.Files
		fmt.Println("Search results:")
		for i, file := range res.Results.Files {
			fmt.Println(i, file.Name)
		}
	}
	return false
}

func (nd *NodeInfo) ExposeFiles() error {
	var files []*t.File
	err := filepath.Walk(nd.localpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			FileNameSplit := strings.Split(info.Name(), ".")
			FILE := &t.File{
				Id:       int32(info.ModTime().Unix()),
				Name:     FileNameSplit[0],
				Filename: info.Name(),
				Size:     int32(info.Size()),
				Ownerid:  nd.Id,
				Type:     mime.TypeByExtension("." + FileNameSplit[len(FileNameSplit)-1]),
			}
			files = append(files, FILE)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error exposing files:", err)
		return err
	}
	_, err = nd.ctx.GetPeerFiles(context.Background(), &t.FileList{
		Files: files,
	})
	return err

}

func (nd *NodeInfo) GetFile(index string) bool {
	ind, err := strconv.Atoi(index)
	if err != nil {
		fmt.Println("Invalid file index:", err)
		return false
	}
	if len(nd.SearchedFiles) == 0 && ind >= len(nd.SearchedFiles) {
		fmt.Println("No file found at index:", ind)
		if err := nd.ExposeFiles(); err != nil {
			fmt.Println("Error exposing files:", err)
		}
		return false
	}
	file := nd.SearchedFiles[ind]
	FileName, TargetId := file.Filename, file.Ownerid
	fmt.Println("Downloading file:", FileName)
	res, err := nd.ctx.GetPeerConnexion(context.Background(), &s.PeerId{
		Id: TargetId,
	})
	if err != nil {
		fmt.Println("Error getting file:", err)
		return strings.Contains(err.Error(), "Peer not registered")
	}
	if err := nd.DownloadFile(res.Peeraddress, file.Id, FileName, file.Size); err != nil {
		fmt.Println("Error downloading file:", err)
		return strings.Contains(err.Error(), "Peer not registered")
	}
	return false
}

func (nd *NodeInfo) DownloadFile(add string, fileId int32, fileName string, Size int32) error {
	conn, ConErr := grpc.Dial(add, grpc.WithInsecure())
	if ConErr != nil {
		return ConErr
	}
	defer conn.Close()

	client := p.NewPeerClient(conn)
	res, err := client.SendFile(context.Background(), &t.File{
		Name: fileName,
	})
	if err != nil {
		return err
	}
	defer res.CloseSend()

	f, ferr := os.OpenFile(path.Join(nd.localpath, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if ferr != nil {
		return ferr
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(int64(Size), "Downloading file...")

	for {
		data, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, err := io.Copy(io.MultiWriter(f, bar), bytes.NewReader(data.GetData())); err != nil {
			return err
		}
	}

	// Set the file time to the time of the file on the server
	fileTime := time.Unix(int64(fileId), 0)
	err = os.Chtimes(path.Join(nd.localpath, fileName), fileTime, fileTime)
	if err != nil {
		return err
	}
	return nil
}

func GetRandomPort() (string, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	defer lis.Close()
	return strings.Split(lis.Addr().String(), ":")[3], nil
}

func InputReader(r *bufio.Reader) (string, error) {
	text, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Replace(text, "\n", "", -1), nil
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (nd *NodeInfo) StartPeerServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := p.Peer{Path: nd.localpath}
	nd.grpcServer = grpc.NewServer()
	p.RegisterPeerServer(nd.grpcServer, &s)

	if err := nd.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (nd *NodeInfo) StopPeerClient() {
	if nd.grpcServer != nil {
		nd.grpcServer.Stop()
	}
}
