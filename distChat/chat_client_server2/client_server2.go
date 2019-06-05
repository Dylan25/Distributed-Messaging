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
	"google.golang.org/grpc/credentials"
)

var (
	collection *mongo.Collection
	crt        = "../authentication/server.crt"
	key        = "../authentication/server.pem"
)

type server struct {
	Broadcast chan chatpb.ChatResponse
	Listener  chan chatpb.ListenResponse

	Users       map[string]string
	UserStreams map[string]chan chatpb.ChatResponse

	StreamMutex sync.RWMutex
	UserMutex   sync.RWMutex
}

func (c *server) Chat(stream chatpb.ChatService_ChatServer) error {
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
	fmt.Printf("searching for messages in: %v\n", name)
	history := database.GetAllMessagesInGroup(name, collection)
	//fmt.Printf("history at 0 is: %v\n", history[0])

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
	fmt.Println("server doing a Listen")
	c.sendMessages(stream)
	return nil
}

func (c *server) setName(user string) {
	c.UserMutex.Lock()
	c.Users[user] = user
	c.UserMutex.Unlock()
}

func (c *server) delName(user string) {
	c.UserMutex.Lock()
	if c.Users[user] == user {
		delete(c.Users, user)
	}
	c.UserMutex.Unlock()
}

func (c *server) setMessage(user string, req *chatpb.ChatRequest) {
	c.StreamMutex.Lock()

	usr := req.GetMsg().GetUser()
	text := req.GetMsg().GetText()
	time := req.GetMsg().GetTime()
	group := req.GetMsg().GetGroup()

	c.UserStreams[user] <- chatpb.ChatResponse{
		Msg: &chatpb.Letter{
			User:  usr,
			Text:  text,
			Time:  time,
			Group: group,
		},
	}
	c.StreamMutex.Unlock()
}

func (c *server) openStream(user string) {
	stream := make(chan chatpb.ChatResponse, 100)

	c.StreamMutex.Lock()
	c.UserStreams[user] = stream
	c.StreamMutex.Unlock()
}

func (c *server) closeStream(user string) {
	c.StreamMutex.Lock()

	if stream, ok := c.UserStreams[user]; ok {
		delete(c.UserStreams, user)
		close(stream)
	}

	c.StreamMutex.Unlock()
}

func (c *server) broadcast(ctx context.Context) {
	for res := range c.Broadcast {
		c.StreamMutex.RLock()
		for _, stream := range c.UserStreams {
			stream <- res
		}
		c.StreamMutex.RUnlock()
	}
}

func (c *server) sendMessages(srv chatpb.ChatService_ListenServer) {
	req, err := srv.Recv()
	if err != nil {
		log.Fatalf("Error recieving from Listen stream: %v\n", err)
	}
	fmt.Printf("server dis sendMessages recieve group is: %v\n", req.GetGroup())
	c.SendMessageHistory(req.GetGroup())
	fmt.Println("server doing sendMessages")
	stream := c.Listener
	defer close(stream)
	for {
		res := <-stream
		srv.Send(&res)
	}

}

func (c *server) getName(user string) (name string, ok bool) {
	c.UserMutex.Lock()
	name, ok = c.Users[user]
	c.UserMutex.Unlock()
	return name, ok
}

func Run(ip string) {
	fmt.Println("Hello World")
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

	//Create TLS Creds
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		fmt.Println("server tls error")
	}
	fmt.Printf("server creds are: %v", creds)

	s := grpc.NewServer(grpc.Creds(creds))
	chatpb.RegisterChatServiceServer(s, &server{
		Broadcast:   make(chan chatpb.ChatResponse, 1000),
		Users:       make(map[string]string),
		UserStreams: make(map[string]chan chatpb.ChatResponse),
		Listener:    make(chan chatpb.ListenResponse, 1000),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		os.Exit(1)
	}

}
