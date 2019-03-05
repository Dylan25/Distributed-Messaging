package c2server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	chatpb "github.com/Distributed-Messaging/distChat/chatpb"
	"google.golang.org/grpc"
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
			log.Fatalf("Error recieving from chat stream: %v\n", err)
		}
		usr := req.GetMsg().GetUser()
		text := req.GetMsg().GetText()
		c.Listener <- chatpb.ListenResponse{
			Msg: &chatpb.Letter{
				User: usr,
				Text: text,
			},
		}
	}
}

func (c *server) Listen(req *chatpb.ListenRequest, stream chatpb.ChatService_ListenServer) error {
	fmt.Println("server doing listen")
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

	c.UserStreams[user] <- chatpb.ChatResponse{
		Msg: &chatpb.Letter{
			User: usr,
			Text: text,
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
