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
