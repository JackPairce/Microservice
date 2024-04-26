package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	p "github.com/JackPairce/MicroService/services/peer"
	s "github.com/JackPairce/MicroService/services/superpeer"
	t "github.com/JackPairce/MicroService/services/types"
	"google.golang.org/grpc"
)

type NodeInfo struct {
	Name          string
	Pass          string
	Id            int32
	ctx           s.SuperPeerClient
	SearchedFiles []*t.File
	MyServerPort  string
	localpath     string `default:"/home/jackpairce/Documents/"`
}

func (nd *NodeInfo) Register() {
	for {
		res, err := nd.ctx.Register(context.Background(), &s.RegisterRequest{
			Name:     nd.Name,
			Password: nd.Pass,
			Address:  GetLocalIP(),
			Port:     nd.MyServerPort,
		})

		if res.Success {
			fmt.Println("Registration successful! with id: ", res.Id)
			nd.Id = res.Id

			break
		} else {
			fmt.Println("Registration failed:", err)
		}
	}
}
func (nd *NodeInfo) Login() {
	for {
		res, err := nd.ctx.Login(context.Background(), &s.RegisterRequest{
			Name:     nd.Name,
			Password: nd.Pass,
			Address:  GetLocalIP(),
		})

		if res.Success {
			fmt.Println("Login successful!")
			nd.Id = res.Id
			break
		} else {
			fmt.Println("Login failed:", err)
		}
	}
}

func (nd *NodeInfo) SearchFiles(searchTerm string) {
	res, err := nd.ctx.SearchFiles(context.Background(), &s.SearchFilesRequest{
		Id:       nd.Id,
		Filename: searchTerm,
	})

	if err != nil {
		fmt.Println("Error searching files:", err)
		return
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
}

func (nd *NodeInfo) ExposeFiles() error {
	var files []*t.File
	dirPath := path.Join(nd.localpath, nd.Name)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		OwnerID := make([]int32, 0)
		if !info.IsDir() {
			FILE := &t.File{
				Name:    info.Name(),
				Size:    int32(info.Size()),
				Ownerid: append(OwnerID, nd.Id),
				Path:    "/",
			}
			files = append(files, FILE)
		}
		return nil
	})
	nd.ctx.GetPeerFiles(context.Background(), &t.FileList{
		Files: files,
	})
	return err

}

func (nd *NodeInfo) GetFile(index string) {
	ind, err := strconv.Atoi(index)
	if err != nil {
		fmt.Println("Invalid file index:", err)
		return
	}
	if len(nd.SearchedFiles) == 0 && ind >= len(nd.SearchedFiles) {
		fmt.Println("No file found at index:", ind)
		if err := nd.ExposeFiles(); err != nil {
			fmt.Println("Error exposing files:", err)
		}
		return
	}
	file := nd.SearchedFiles[ind]
	FileName, TargetId := file.Name, file.Ownerid[0]
	fmt.Println("Downloading file:", FileName)
	res, err := nd.ctx.GetPeerConnexion(context.Background(), &s.PeerId{
		Id: TargetId,
	})
	if err != nil {
		fmt.Println("Error getting file:", err)
		return
	}
	if err := nd.DownloadFile(res.Address+":"+res.Port, FileName); err != nil {
		fmt.Println("Error downloading file:", err)
	}
}

func (nd *NodeInfo) adjustFileName(fileName *string) {
	index := 1
	splitFile := strings.Split(*fileName, ".")
	baseName := splitFile[0]
	for {
		// Check if the file exists
		_, err := os.Stat(path.Join(nd.localpath, nd.Name, strings.Join(splitFile, ".")))
		if os.IsNotExist(err) {
			break // File does not exist, use this filename
		}
		// File exists, increment index and try again
		splitFile[0] = fmt.Sprintf("%s(%d)", baseName, index)
		index++
	}
	*fileName = strings.Join(splitFile, ".")
}

func (nd *NodeInfo) DownloadFile(add string, fileName string) error {
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

	nd.adjustFileName(&fileName)

	f, ferr := os.Create(path.Join(nd.localpath, nd.Name, fileName))
	if ferr != nil {
		return ferr
	}
	defer f.Close()

	for {
		data, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, err := f.Write(data.GetData()); err != nil {
			return err
		}
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

	s := p.Peer{Path: nd.localpath, Name: nd.Name}
	grpcServer := grpc.NewServer()
	p.RegisterPeerServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
