package server

import (
	"chat/chatpb"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	webPort = 8000
)

type ChatServiceServer struct {
	chatpb.UnimplementedChatServiceServer
	channels map[string][]chan *chatpb.Message
}

func (s *ChatServiceServer) JoinChannel(channel *chatpb.Channel, stream chatpb.ChatService_JoinChannelServer) error {
	msgChannel := make(chan *chatpb.Message)
	s.channels[channel.Name] = append(s.channels[channel.Name], msgChannel)
	log.Printf("Added new channel for '%s' to '%s' channels\n", channel.SendersName, channel.Name)
	log.Printf("number of users in channel '%s': '%d'\n", channel.Name, len(s.channels[channel.Name]))
	for {
		select {
		case <-stream.Context().Done():
			log.Printf("stream context for '%s' on channel '%s' is done, ending the conversation...\n", channel.SendersName, channel.Name)
			return nil
		case msg := <-msgChannel:
			log.Printf("Recieved message from message channel '%s' sent by '%s': %s\n", channel.Name, channel.SendersName, msg)
			stream.Send(msg)
		}
	}

}

func (s *ChatServiceServer) SendMessage(stream chatpb.ChatService_SendMessageServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	if err == io.EOF {
		return nil
	}

	ack := &chatpb.MessageAck{
		Status: "success",
	}
	stream.SendAndClose(ack)
	channels := s.channels[msg.Channel.Name]
	defer func() {
		if x := recover(); x != nil {
			log.Printf("unable to send: %v\n", x)
		}
	}()
	for _, msgChan := range channels {
		log.Printf("sending message '%s' to users in channel...\n", msg.Message)
		msgChan <- msg
		log.Println("done")
	}
	return nil
}

func Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatal(err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chatpb.RegisterChatServiceServer(grpcServer, &ChatServiceServer{
		channels: make(map[string][]chan *chatpb.Message),
	})

	log.Println("Server is listening on port:", webPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
