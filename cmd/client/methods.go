package main

import (
	"context"
	"fmt"

	"github.com/JackPairce/MicroService/services/superpeer"
)

func Register(c superpeer.SuperPeerClient, Name string, password string) {
	for {
		res, err := c.Register(context.Background(), &superpeer.RegisterRequest{
			Name:     Name,
			Password: password,
			Address:  GetLocalIP(),
		})

		if res.Success {
			fmt.Println("Registration successful!")
			break
		} else {
			fmt.Println("Registration failed:", err)
		}
	}
}
func Login(c superpeer.SuperPeerClient, Name string, password string) {
	for {
		res, err := c.Login(context.Background(), &superpeer.RegisterRequest{
			Name:     Name,
			Password: password,
			Address:  GetLocalIP(),
		})

		if res.Success {
			fmt.Println("Login successful!")
			break
		} else {
			fmt.Println("Login failed:", err)
		}
	}
}

func SearchFiles(c superpeer.SuperPeerClient, searchTerm string) {
	res, err := c.SearchFiles(context.Background(), &superpeer.SearchFilesRequest{
		Id:       0,
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
