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
	//chatting constructs and sends a users message to a single target
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

func chatConsole(clients []chatpb.ChatServiceClient, group string, user string, l chatpb.ChatServiceClient) {
	//ChatConsole uses two go routines to send the users messages, and recieve messages for the user.
	waitc := make(chan struct{})
	exitChannel := make(chan struct{})
	listenStopped := make(chan struct{})

	//create the listener stream
	stream, err := l.Listen(context.Background())
	if err != nil {
		log.Fatalf("Error creating Stream: %v", err)
		return
	}

	//send outgoing messages
	go func() {
		fmt.Printf("Joined %v as %v. Type '!exit' to return to the menu\n", group, user)

		buf := bufio.NewReader(os.Stdin)
		for {
			text, err := buf.ReadString('\n')
			if strings.ToLower(text) == "!exit\n" {
				//Leave the chatroom, shutdown both go routines
				close(exitChannel)
				//send a message to the listen server so that
				//its go routine switch statment can see that exitChannel has been closed
				stream.Send(&chatpb.ListenRequest{
					User:  user,
					Group: group,
				})
				//wait for the Listen go routine to stop
				<-listenStopped
				//exit
				close(waitc)
				return
			}
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

	//Listen for incoming messages
	go func() {
		stream.Send(&chatpb.ListenRequest{
			User:  user,
			Group: group,
		})

		for {
			select {
			default:
				res, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("Erro while recieving: %v", err)
					break
				}
				//display message to user
				user := res.GetMsg().GetUser()
				messagetime := time.Unix(res.GetMsg().GetTime(), 0)
				text := res.GetMsg().GetText()
				text = text[:len(text)-1]
				reply := ":" + user + ": " + text
				fmt.Printf("%v%s\n", messagetime, reply)
				//the below case closes this function
			case <-exitChannel:
				close(listenStopped)
				return
			}
		}
	}()
	<-waitc
}

func makeClients(Ips []string) []chatpb.ChatServiceClient {
	//makeClients generates a Grpc client stub for each ip in the group
	//the client is connecting to
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
	//RunMenu maintains the user's account state sets up the listener server and
	//drives the menu parser
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
	//ParseMenuInput calls the apropriate menu function on the users input
	switch input {
	case "!help":
		fmt.Println("!createaccount, !signIn, !changepassword, !joingroup, !listgroups, !help")
	case "!createaccount":
		menu.CreateAccount()
	case "!signin":
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
			fmt.Printf("connectint to: %v at %v", groupName, ips)
			if err != nil {
				fmt.Printf("JoinGroup error: %v", err)
			}
			clients := makeClients(ips)
			chatConsole(clients, groupName, usersAccount.Name, l)
			fmt.Println("done chatting")
		} else {
			fmt.Println("Please signin before joining a group")
		}
	case "!listgroups":
		menu.ListGroups(usersAccount.Name, *signedIn)
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
			fmt.Printf("Could not connect to: %v error: %v\n", ip, lerr)
			port++
			continue
		}
		testIfOpen.Close()
		//found an open port so run server on it
		fmt.Printf("Hosting on: %v\n", ip)
		go c2server.Run(ip)
		lc, lerr = grpc.Dial(ip, grpc.WithInsecure())
		if lerr != nil {
			port++
		} else {
			break
		}
	}
	return lc
}

func main() {

	//connect to mongo
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
	RunMenu()
}
