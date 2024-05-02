package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/JackPairce/MicroService/services/superpeer"
	t "github.com/JackPairce/MicroService/services/types"
)

const (
	port    = "8080"
	address = "localhost"
)

var (
	MyID int64
)

func main() {
	conn, err := grpc.Dial(address+":"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			func(ctx context.Context,
				method string,
				req,
				reply interface{},
				cc *grpc.ClientConn,
				invoker grpc.UnaryInvoker,
				opts ...grpc.CallOption,
			) error {
				md, ok := metadata.FromOutgoingContext(ctx)
				if ok {
					log.Println("-->out md interceptor: ", md)
				}
				return invoker(ctx, method, req, reply, cc, opts...)
			}),
	)
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

	nd := NodeInfo{
		ctx:           c,
		SearchedFiles: []*t.File{},
		MyServerPort:  MyPort,
	}

Choose:
	for {
		fmt.Print("chose option:\n1. Register\n2. Login\n->")
		option, _ := InputReader(reader)
		if option != "1" && option != "2" {
			fmt.Println("Invalid option")
			continue
		}

		fmt.Print("Enter your Name: ")
		Name, _ := InputReader(reader)
		fmt.Print("Enter your PassWord: ")
		password, _ := InputReader(reader)

		nd.Name = Name
		nd.Pass = password
		switch option {
		case "1":
			if nd.Register() {
				break Choose
			}
		case "2":
			if nd.Login() {
				break Choose
			}
		}
	}
	go nd.StartPeerServer(MyPort)

	// TODO: Add a notifier to add path
	// TODO: Add a watch on the path for changes

	nd.localpath = path.Join("/home/jackpairce/Documents/", nd.Name)
	fmt.Println("Exposing Files from:", nd.localpath)
	if err := os.Mkdir(nd.localpath, 0755); err != nil {
		if !os.IsExist(err) {
			log.Fatalln(err)
		}
	}

	if err := nd.ExposeFiles(); err != nil {
		log.Fatalln(err)
	}

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
		if xComand[0] == "/exit" {
			return
		}
		if len(xComand) < 2 {
			fmt.Println("Invalid Command")
			continue
		}
		C := xComand[0]
		arg := xComand[1]
		switch C {
		case "":
		case "/find":
			if nd.SearchFiles(arg) {
				nd.StopPeerClient()
				goto Choose
			}
		case "/get":
			if nd.GetFile(arg) {
				nd.StopPeerClient()
				goto Choose
			}
		default:
			fmt.Println("Command Not reconized of (/find,/get,/exit)")
		}
	}
}

// Unary client interceptor to add custom headers to outgoing requests
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	// Add custom headers to outgoing requests
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("id", strconv.Itoa(int(123))))
	// Call the RPC invoker
	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}
