package client

import (
	chat "chat/proto"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

const webPort = 8000

func Run() {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", webPort), grpc.WithInsecure())
	if err != nil {
		log.Fatal("error connecting to server on port:", webPort)
	}

	client := chat.NewChatServiceClient(conn)
	msg := &chat.Message{
		From: "client",
		Body: "This is my message!",
	}

	steam, err := client.Join(context.Background())
}
