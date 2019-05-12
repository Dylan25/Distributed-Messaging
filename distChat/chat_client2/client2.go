package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	c2server "github.com/Distributed-Messaging/distChat/chat_client_server2"
	chatpb "github.com/Distributed-Messaging/distChat/chatpb"
	database "github.com/Distributed-Messaging/distChat/database"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

func chatting(c chatpb.ChatServiceClient, user string, text string, time int64, group string) {
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("Error creating Stream: %v", err)
		return
	}

	stream.Send(&chatpb.ChatRequest{
		Msg: &chatpb.Letter{
			User:  user,
			Text:  text,
			Time:  time,
			Group: group,
		},
	})
}

func chatConsole(clients []chatpb.ChatServiceClient, group string, l chatpb.ChatServiceClient) {
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
		go listening(l, user, group)

		for {
			text, err := buf.ReadString('\n')
			time := int64(time.Now().Unix())
			if err != nil {
				log.Fatalf("Error reading message input: %v", err)
				continue
			}

			for _, c := range clients {
				chatting(c, user, text, time, group)
			}
		}
	}()

	<-waitc

}

func listening(c chatpb.ChatServiceClient, user string, group string) {
	stream, err := c.Listen(context.Background())
	if err != nil {
		log.Fatalf("Error creating Stream: %v", err)
		return
	}
	stream.Send(&chatpb.ListenRequest{
		User:  user,
		Group: group,
	})

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
		messagetime := time.Unix(res.GetMsg().GetTime(), 0)
		user = user[:len(user)-1]
		text := res.GetMsg().GetText()
		text = text[:len(text)-1]
		reply := ":" + user + ": " + text
		fmt.Printf("%v%s\n", messagetime, reply)
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

func pickGroup() ([]string, string) {

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("enter group name: ")
	group, grouperr := buf.ReadString('\n')
	if grouperr != nil {
		log.Fatalf("Error reading ip and port: %v", grouperr)
	}
	group = group[:len(group)-1]

	grouptojoin, err := database.GetOneGroup(group, collection)

	if err != nil {
		fmt.Print("ips to connect to: ")
		ipToConnect, ipcerr := buf.ReadString('\n')
		if ipcerr != nil {
			log.Fatalf("Error reading ip and port: %v", ipcerr)
		}
		ipToConnect = ipToConnect[:len(ipToConnect)-1]
		IPs := strings.Fields(ipToConnect)

		grouptojoin = database.Group{
			IPs:  IPs,
			Name: group,
		}

		database.StoreGroup(grouptojoin, collection)
	}

	fmt.Printf("connecting to %s", group)
	return grouptojoin.IPs, grouptojoin.Name
}

func main() {

	uri := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("mydb").Collection("chatgroups")

	// fmt.Println("group: %v", database.GetOneGroup("friends", collection))

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("server ip and port: ")
	ip, iperr := buf.ReadString('\n')
	if iperr != nil {
		log.Fatalf("Error reading ip and port: %v", iperr)
	}
	ip = ip[:len(ip)-1]
	go c2server.Run(ip)

	// fmt.Print("ips to connect to: ")
	// ipToConnect, ipcerr := buf.ReadString('\n')
	// if ipcerr != nil {
	// 	log.Fatalf("Error reading ip and port: %v", ipcerr)
	// }
	// ipToConnect = ipToConnect[:len(ipToConnect)-1]
	// IPs := strings.Fields(ipToConnect)
	// fmt.Println(IPs)

	IPs, groupName := pickGroup()

	clients := makeClients(IPs)

	//c := chatpb.NewChatServiceClient(cc)

	lc, lerr := grpc.Dial(ip, grpc.WithInsecure())
	if lerr != nil {
		log.Fatalf("Could not connect: %v", lerr)
	}
	defer lc.Close()

	l := chatpb.NewChatServiceClient(lc)

	//fmt.Println("prepairing to chat\n")

	chatConsole(clients, groupName, l)

}
