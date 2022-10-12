package server

import (
	chat "chat/proto"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var messageChannel = make(chan chat.Message, 1)

const webPort = 8000

type Server struct{}

func (s Server) Join(emp *empty.Empty, server chat.ChatService_JoinServer) error {
	msg := <-messageChannel
	if err := server.Send(&msg); err != nil {
		log.Println("Error sending message.")
	}
	return nil
}

func (s Server) SendMessage(ctx context.Context, message *chat.Message) (*empty.Empty, error) {
	messageChannel <- *message
	return &empty.Empty{}, nil
}

func Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatal("error listening on port:", webPort)
	}

	grpcServer := grpc.NewServer()
	server := Server{}

	chat.RegisterChatServiceServer(grpcServer, server)

	log.Println("Starting server on port:", webPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("error running server on port:", webPort)
	}
}
