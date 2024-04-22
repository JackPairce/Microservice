package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	s "github.com/JackPairce/MicroService/services/superpeer"
	t "github.com/JackPairce/MicroService/services/types"
)

type NodeInfo struct {
	Name string
	Pass string
	Id   int32
	ctx  s.SuperPeerClient
}

func (nd *NodeInfo) Register() {
	for {
		res, err := nd.ctx.Register(context.Background(), &s.RegisterRequest{
			Name:     nd.Name,
			Password: nd.Pass,
			Address:  GetLocalIP(),
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
		fmt.Println("Search results:")
		for _, file := range res.Results.Files {
			fmt.Println("- ", file.Name)
		}
	}
}

func (nd *NodeInfo) ExposeFiles() error {
	var files []*t.File
	dirPath := "/home/jackpairce/Documents/" + nd.Name
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
