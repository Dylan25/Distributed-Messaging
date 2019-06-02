package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	menu "github.com/Distributed-Messaging/distChat/chat_client2_menu"
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

func RunMenu() {
	var signedIn bool
	var usersAccount database.Account
	fmt.Println("Welcome to DistChat")
	fmt.Println("type '!help' for a list of commands")
	lc := GetListeningConnection()
	defer lc.Close()
	l := chatpb.NewChatServiceClient(lc)
	for {
		buf := bufio.NewReader(os.Stdin)
		input, inputerr := buf.ReadString('\n')
		if inputerr != nil {
			fmt.Printf("Error reading password: %v", inputerr)
		}
		input = input[:len(input)-1]
		ParseMenuInput(strings.ToLower(input), &signedIn, &usersAccount, l)
	}
}

func ParseMenuInput(input string, signedIn *bool, usersAccount *database.Account, l chatpb.ChatServiceClient) {
	switch input {
	case "!help":
		fmt.Println("!createaccount, !signIn, !changepassword, !joingroup, !help")
	case "!createaccount":
		menu.CreateAccount()
	case "!signin":
		//var account database.Account
		var signInError error
		*signedIn, *usersAccount, signInError = menu.SignIn()
		*signedIn = false
		if signInError != nil {
			fmt.Printf("Sign in Error: %v\n", signInError)
		} else {
			*signedIn = true
		}
	case "!changepassword":
		menu.ChangePassword()
	case "!joingroup":
		if *signedIn {
			ips, groupName, err := menu.PickGroup(usersAccount.Name)
			if err != nil {
				fmt.Printf("JoinGroup error: %v", err)
			}
			clients := makeClients(ips)
			chatConsole(clients, groupName, l)
		} else {
			fmt.Println("Please signin before joining a group")
		}
	}
}

func GetListeningConnection() *grpc.ClientConn {
	//GetListeningConnection finds an open port to host the listener server on.
	var lc *grpc.ClientConn
	var lerr error
	port := 50051
	for {
		//find an open port
		ip := "0.0.0.0:" + strconv.Itoa(port)
		testIfOpen, testError := net.Listen("tcp", ":"+strconv.Itoa(port))
		if testError != nil {
			fmt.Printf("1Could not connect to: %v error: %v", ip, lerr)
			port++
			continue
		}
		testIfOpen.Close()
		//found an open port so run server on it
		go c2server.Run(ip)
		lc, lerr = grpc.Dial(ip, grpc.WithInsecure())
		if lerr != nil {
			fmt.Printf("2Could not connect to: %v error: %v", ip, lerr)
			port++
		} else {
			break
		}
	}
	return lc
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
	menu.SetCollection(collection)

	// fmt.Println("group: %v", database.GetOneGroup("friends", collection))

	// buf := bufio.NewReader(os.Stdin)
	// fmt.Print("server ip and port: ")
	// ip, iperr := buf.ReadString('\n')
	// if iperr != nil {
	// 	log.Fatalf("Error reading ip and port: %v", iperr)
	// }
	// ip = ip[:len(ip)-1]
	// go c2server.Run(ip)
	// IPs, groupName := pickGroup()
	// clients := makeClients(IPs)
	// lc, lerr := grpc.Dial(ip, grpc.WithInsecure())
	// if lerr != nil {
	// 	log.Fatalf("Could not connect: %v", lerr)
	// }
	// defer lc.Close()
	// l := chatpb.NewChatServiceClient(lc)
	// chatConsole(clients, groupName, l)

	RunMenu()

}
