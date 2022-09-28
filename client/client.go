package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xvbnm48/grpc-chat/chatserver/chatserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("enter server port ::")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read port server: %v", err)
	}

	serverID = strings.Trim(serverID, "\r\n")

	log.Println("connecting: " + serverID)

	conn, err := grpc.Dial(serverID, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connection to server: %v", err)
	}
	defer conn.Close()
	// call chat service
	client := chatserver.NewServicesClient(conn)
	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("failed to call chatservice:  %v", err)
	}

	// implement connection to server
	ch := clientHandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()
	// blocker
	bl := make(chan bool)
	<-bl
}

type clientHandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("enter your name :: ")
	Name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read from console :: %v", err)
	}

	ch.clientName = strings.Trim(Name, "\r\n")
}

// send message
func (ch *clientHandle) sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)
		ClientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read from console ::%v ", err)
		}

		ClientMessage = strings.Trim(ClientMessage, "\r\n")
		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: ClientMessage,
		}

		err = ch.stream.Send(clientMessageBox)
		if err != nil {
			log.Printf("error while sending message: %v ", err)
		}
	}
}

func (ch *clientHandle) receiveMessage() {
	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}
		fmt.Printf("%s : %s \n", mssg.Name, mssg.Body)
	}
}
