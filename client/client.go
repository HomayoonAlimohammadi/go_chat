package client

import (
	"bufio"
	"chat/chatpb"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const webPort = 8000

// server-side streaming client
func joinChannel(ctx context.Context, client chatpb.ChatServiceClient, channelName, sendersName string) {
	channel := chatpb.Channel{Name: channelName, SendersName: sendersName}
	stream, err := client.JoinChannel(ctx, &channel)
	if err != nil {
		log.Fatal("can not join channel:", channelName)
	}

	log.Println("joined channel:", channelName)

	// waitChannel := make(chan struct{})

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Printf("channel '%s' is closing", channelName)
			return
		}
		if err != nil {
			log.Fatal("error receiving message from channel:", channelName, err)
		}

		log.Printf("%s: %s\n", msg.Sender, msg.Message)
	}
	// <-waitChannel
}

// client-side streaming client
func sendMessage(ctx context.Context, client chatpb.ChatServiceClient, messageText, channelName, sendersName string) error {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Fatal("can not send message:", err)
	}
	channel := &chatpb.Channel{
		Name:        channelName,
		SendersName: sendersName,
	}
	message := &chatpb.Message{
		Sender:  sendersName,
		Channel: channel,
		Message: messageText,
	}
	// log.Printf("sending message to channel '%s': %s\n", channelName, message)
	err = stream.Send(message)
	if err != nil {
		log.Fatal("error sending message to all users...")
	}
	ack, err := stream.CloseAndRecv()
	if err != nil || ack.Status != "success" {
		log.Println("can not receive message acknowledgement from the server:", err)
	}
	return nil
}

func Run(channelName, sendersName string, interval int) {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", webPort), opts...)
	if err != nil {
		log.Fatal("error connecting to server on port:", webPort)
	}
	defer conn.Close()

	ctx := context.Background()
	client := chatpb.NewChatServiceClient(conn)
	go joinChannel(ctx, client, channelName, sendersName)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text(), channelName, sendersName)
	}
}
