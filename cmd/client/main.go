package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"

	"github.com/JackPairce/MicroService/services/superpeer"
)

func main() {
	port := "8080"
	address := "localhost"
	// Set up a connection to the server.
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := superpeer.NewSuperPeerClient(conn)
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
		ctx:  c,
		Name: Name,
		Pass: password,
	}
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
		fmt.Print("Enter File Name To find: ")
		FileToSearch, err := InputReader(reader)
		if err != nil {
			log.Fatalln(err)
			return
		}
		if FileToSearch == "/exit" {
			break
		}
		nd.SearchFiles(FileToSearch)

	}
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
