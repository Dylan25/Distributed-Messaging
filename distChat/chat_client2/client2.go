package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	c2server "github.com/Distributed-Messaging/distChat/chat_client_server2"
	chatpb "github.com/Distributed-Messaging/distChat/chatpb"
	"google.golang.org/grpc"
)

func chatting(c chatpb.ChatServiceClient, user string, text string) {
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("Error creating Stream: %v", err)
		return
	}

	stream.Send(&chatpb.ChatRequest{
		Msg: &chatpb.Letter{
			User: user,
			Text: text,
		},
	})
}

func chatConsole(clients []chatpb.ChatServiceClient) {
	waitc := make(chan struct{})
	go func() {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("Enter user name: ")
		user := ""
		var uerr error
		for user == "" {
			user, uerr = buf.ReadString('\n')
			if uerr != nil {
				log.Fatalf("Error reading message input: %v", uerr)
			}
		}

		for {
			text, err := buf.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading message input: %v", err)
				continue
			}

			for _, c := range clients {
				chatting(c, user, text)
			}
		}
	}()

	<-waitc

}

func listening(c chatpb.ChatServiceClient) {
	stream, err := c.Listen(context.Background(), &chatpb.ListenRequest{
		User: "owner",
	})
	if err != nil {
		log.Fatalf("Error creating Stream: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Erro while recieving: %v", err)
			break
		}
		user := res.GetMsg().GetUser()
		user = user[:len(user)-1]
		text := res.GetMsg().GetText()
		text = text[:len(text)-1]
		reply := user + ": " + text
		fmt.Println(reply)
	}

}

func makeClients(Ips []string) []chatpb.ChatServiceClient {
	var clients []chatpb.ChatServiceClient
	for _, ip := range Ips {
		cc, err := grpc.Dial(ip, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %v", err)
		}
		c := chatpb.NewChatServiceClient(cc)
		clients = append(clients, c)
	}

	return clients
}

func main() {

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("server ip and port: ")
	ip, iperr := buf.ReadString('\n')
	if iperr != nil {
		log.Fatalf("Error reading ip and port: %v", iperr)
	}
	ip = ip[:len(ip)-1]
	go c2server.Run(ip)

	fmt.Print("ips to connect to: ")
	ipToConnect, ipcerr := buf.ReadString('\n')
	if ipcerr != nil {
		log.Fatalf("Error reading ip and port: %v", ipcerr)
	}
	ipToConnect = ipToConnect[:len(ipToConnect)-1]
	IPs := strings.Fields(ipToConnect)
	fmt.Println(IPs)

	clients := makeClients(IPs)

	//c := chatpb.NewChatServiceClient(cc)

	lc, lerr := grpc.Dial(ip, grpc.WithInsecure())
	if lerr != nil {
		log.Fatalf("Could not connect: %v", lerr)
	}
	defer lc.Close()

	l := chatpb.NewChatServiceClient(lc)

	go listening(l)
	//fmt.Println("prepairing to chat\n")
	chatConsole(clients)
}
