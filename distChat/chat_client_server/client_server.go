package cserver

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	chatpb "github.com/dylChat/chatpb"
	"google.golang.org/grpc"
)

type server struct {
	Broadcast chan chatpb.ChatResponse

	Users       map[string]string
	UserStreams map[string]chan chatpb.ChatResponse

	StreamMutex sync.RWMutex
	UserMutex   sync.RWMutex
}

// func (*server) Connect(context.Background(), req *chatpb.ConnectRequest) (res *chatpb.ConnectResponse, error) {
// 	ip := req.GetIp()

// }

//March 15th

func (c *server) Chat(stream chatpb.ChatService_ChatServer) error {
	//go c.sendBroadcasts(stream, users)

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

		_, ok := c.getName(usr)
		if ok == false {
			c.setName(usr)
			go c.sendMessages(stream, usr)
			//fmt.Printf("set: %v\n", usr)
			go c.broadcast(context.Background())
			//fmt.Println("after send Messages")
			//c.setMessage(usr, req)
			//fmt.Println("after set message")
		}

		c.Broadcast <- chatpb.ChatResponse{
			Msg: &chatpb.Letter{
				User: usr,
				Text: text,
			},
		}

		// for _, users := range c.Users {
		// 	c.UserStreams[users].Send(chatpb.ChatResponse{
		// 		Msg: &chatpb.Letter{
		// 			User: usr,
		// 			Text: text,
		// 		},
		// 	})
		// }

		// c.Broadcast <- chatpb.ChatResponse{
		// 	Msg: &chatpb.Letter{
		// 		User: usr,
		// 		Text: text,
		// 	},
		// }

		//c.broadcast(context.Background())

		// if serr != nil {
		// 	log.Fatalf("Error sending msg to clients: %v\n", serr)
		// 	return serr
		// }
	}
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
	//fmt.Println("in setMessage")
	c.StreamMutex.Lock()
	//fmt.Println("in setMessage with lock")

	usr := req.GetMsg().GetUser()
	text := req.GetMsg().GetText()

	//fmt.Println("before message added")
	c.UserStreams[user] <- chatpb.ChatResponse{
		Msg: &chatpb.Letter{
			User: usr,
			Text: text,
		},
	}
	//fmt.Println("after message added")

	//fmt.Printf("UserStreams at %v is: %v", user, <-c.UserStreams[user])
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

	//DebugLogf("closed stream for client %s", user)

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

func (c *server) sendMessages(srv chatpb.ChatService_ChatServer, user string) {
	c.openStream(user)
	stream := c.UserStreams[user]
	defer c.closeStream(user)

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

// func (s *server) sendBroadcasts(srv chatpb.ChatService_ChatServer, user string) {
// 	stream := s.openStream(user)
// 	defer s.closeStream(user)

// 	for {
// 		select {
// 		case <-srv.Context().Done():
// 			return
// 		case res := <-stream:
// 			srv.Send(&res)
// 		}
// 	} //TODO refactor this
// }

//TODO deal with disconnecting and dropping clients etc.

//func Run()

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
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		os.Exit(1)
	}

}
