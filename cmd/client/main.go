package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"

	"github.com/JackPairce/MicroService/services/superpeer"
	t "github.com/JackPairce/MicroService/services/types"
)

const (
	port    = "8080"
	address = "localhost"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := superpeer.NewSuperPeerClient(conn)
	var MyPort string
	MyPort, err = GetRandomPort()
	if err != nil {
		log.Fatalln(err)
	}
	reader := bufio.NewReader(os.Stdin)
	option := "0"
	for {
		fmt.Print("chose option:\n1. Register\n2. Login\n->")
		option, _ = InputReader(reader)
		if option == "1" || option == "2" {
			break
		} else {
			fmt.Println("Invalid option")
		}
	}

	fmt.Print("Enter your Name: ")
	Name, _ := InputReader(reader)
	fmt.Print("Enter your PassWord: ")
	password, _ := InputReader(reader)
	nd := NodeInfo{
		ctx:           c,
		Name:          Name,
		Pass:          password,
		SearchedFiles: []*t.File{},
		MyServerPort:  MyPort,
		localpath:     "/home/jackpairce/Documents/",
	}
	go nd.StartPeerServer(MyPort)

	switch option {
	case "1":
		nd.Register()
	case "2":
		nd.Login()
	default:
		fmt.Println("Invalid option")
	}
	nd.ExposeFiles()

	for {
		fmt.Print("Enter Command (/find,/get,/exit)\n> ")
		Command, err := InputReader(reader)
		if err != nil {
			log.Fatalln(err)
			return
		}
		if Command == "" {
			continue
		}
		xComand := strings.Split(Command, " ")
		if len(xComand) < 2 {
			fmt.Println("Invalid Command")
			continue
		}
		C := xComand[0]
		arg := xComand[1]
		switch C {
		case "":
		case "/find":
			nd.SearchFiles(arg)
		case "/get":
			nd.GetFile(arg)
		case "/exit":
			return
		default:
			fmt.Println("Command Not reconized of (/find,/get,/exit)")
		}
	}
}
