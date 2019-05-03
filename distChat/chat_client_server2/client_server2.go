package c2server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	chatpb "github.com/Distributed-Messaging/distChat/chatpb"
	"google.golang.org/grpc"
)

//********* database stuff ********

var collection *mongo.Collection

type Message struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	User string             `bson:"user,omitempty"`
	Text string             `bson:"text,omitempty"`
	Time int64              `bson:"time,omitempty"`
}

func StoreMessage(req *chatpb.ChatRequest) {
	message := req.GetMsg()
	messagetostore := Message{
		User: message.GetUser(),
		Text: message.GetText(),
		Time: message.GetTime(),
	}

	collection.InsertOne(context.Background(), messagetostore)

	//return NOTE do I need the status errors?
}

func GetAllMessages() []Message {
	var messages []Message
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("could not read messages")
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var message Message
		cursor.Decode(&message)
		messages = append(messages, message)
	}

	return messages
}

//********** end database stuff **********

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
			log.Fatalf("Error recieving from chat stream: %v\n", err)
		}
		StoreMessage(req)
		usr := req.GetMsg().GetUser()
		text := req.GetMsg().GetText()
		time := req.GetMsg().GetTime()
		c.Listener <- chatpb.ListenResponse{
			Msg: &chatpb.Letter{
				User: usr,
				Text: text,
				Time: time,
			},
		}
	}
}

func (c *server) SendMessageHistory() {
	history := GetAllMessages()

	for _, message := range history {
		c.Listener <- chatpb.ListenResponse{
			Msg: &chatpb.Letter{
				User: message.User,
				Text: message.Text,
				Time: message.Time,
			},
		}
	}
}

func (c *server) Listen(stream chatpb.ChatService_ListenServer) error {
	fmt.Println("server doing listen")
	c.SendMessageHistory()
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

	c.UserStreams[user] <- chatpb.ChatResponse{
		Msg: &chatpb.Letter{
			User: usr,
			Text: text,
			Time: time,
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

	s := grpc.NewServer()
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
