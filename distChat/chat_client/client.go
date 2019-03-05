package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	cserver "github.com/Distributed-Messaging/distChat/chat_client_server"
	chatpb "github.com/Distributed-Messaging/distChat/chatpb"
	"google.golang.org/grpc"
)

func chatting(c chatpb.ChatServiceClient) {
	waitc := make(chan struct{})
	// messages := make(chan string)
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("Error creating Stream: %v", err)
		return
	}

	//var replies []string
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
			// for _, message := range replies {
			// 	fmt.Println(message)
			// }
			// replies = replies[:0]
			// fmt.Print(":> ")
			text, err := buf.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading message input: %v", err)
				continue
			}
			stream.Send(&chatpb.ChatRequest{
				Msg: &chatpb.Letter{
					User: user,
					Text: text,
				},
			})
		}
	}()

	go func() {
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
			// append(replies, reply) //Race condition?
		}
	}()
	<-waitc

}

func main() {

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("server ip and port: ")
	ip, iperr := buf.ReadString('\n')
	if iperr != nil {
		log.Fatalf("Error reading ip and port: %v", iperr)
	}
	ip = ip[:len(ip)-1]
	go cserver.Run(ip)

	fmt.Print("ip to connect to: ")
	ipToConnect, ipcerr := buf.ReadString('\n')
	if ipcerr != nil {
		log.Fatalf("Error reading ip and port: %v", ipcerr)
	}
	ipToConnect = ipToConnect[:len(ipToConnect)-1]

	fmt.Println("Hello I'm client")
	cc, err := grpc.Dial(ipToConnect, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := chatpb.NewChatServiceClient(cc)

	chatting(c)
	//cs := greetpb.NewSumServiceClient //TODO move sum into its own server client thing
	// fmt.Printf("Created client: %f", c)

	//doUnary(c)

	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
}
