package c2server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	chatpb "github.com/Distributed-Messaging/distChat/chatpb"
	database "github.com/Distributed-Messaging/distChat/database"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct {
	Listener chan chatpb.ListenResponse

	Users       map[string]string
	UserStreams map[string]chan chatpb.ChatResponse

	StreamMutex sync.RWMutex
	UserMutex   sync.RWMutex
}

func (c *server) Chat(stream chatpb.ChatService_ChatServer) error {
	//recieves incoming messages and queues them for later dispatch to the recipient
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("Error recieving from chat stream: %v\n", err)
		}
		database.StoreMessage(req, collection)
		usr := req.GetMsg().GetUser()
		text := req.GetMsg().GetText()
		time := req.GetMsg().GetTime()
		group := req.GetMsg().GetGroup()
		c.Listener <- chatpb.ListenResponse{
			Msg: &chatpb.Letter{
				User:  usr,
				Text:  text,
				Time:  time,
				Group: group,
			},
		}
	}
}

func (c *server) SendMessageHistory(name string) {
	//queues all messages in a group's history for the recipient
	fmt.Printf("searching for messages in: %v\n", name)
	history := database.GetAllMessagesInGroup(name, collection)

	for _, message := range history {
		c.Listener <- chatpb.ListenResponse{
			Msg: &chatpb.Letter{
				User:  message.User,
				Text:  message.Text,
				Time:  message.Time,
				Group: message.Group,
			},
		}
	}
}

func (c *server) Listen(stream chatpb.ChatService_ListenServer) error {
	c.sendMessages(stream)
	return nil
}

func (c *server) sendMessages(srv chatpb.ChatService_ListenServer) {
	//send all queued historical and recieved messages to the recipient
	req, err := srv.Recv()
	go func() {
		for {
			req, err = srv.Recv()
			c.Listener <- chatpb.ListenResponse{
				Msg: &chatpb.Letter{
					User:  req.GetUser(),
					Text:  "leaving group chat.",
					Time:  0,
					Group: req.GetGroup(),
				}}
			return
		}
	}()
	fmt.Printf("server is sendMessages recieve group is: %v\n", req.GetGroup())
	c.SendMessageHistory(req.GetGroup())
	fmt.Println("server doing sendMessages")
	stream := c.Listener
	defer close(stream)
	for {
		res := <-stream
		srv.Send(&res)
	}

}

func Run(ip string) {
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

	collection = client.Database("mydb").Collection("chat")

	lis, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//initialize server
	s := grpc.NewServer()
	chatpb.RegisterChatServiceServer(s, &server{
		Users:       make(map[string]string),
		UserStreams: make(map[string]chan chatpb.ChatResponse),
		Listener:    make(chan chatpb.ListenResponse, 1000),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		os.Exit(1)
	}

}
