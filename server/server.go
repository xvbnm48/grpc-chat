package main

import (
	"log"
	"net"
	"os"

	"github.com/xvbnm48/grpc-chat/chatserver/chatserver"
	"google.golang.org/grpc"
)

func main() {
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000"
	}

	// init listener
	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("could not listen @ %v :: %v ", Port, err)
	}

	log.Println("Listening on port", Port)

	//grpc server instance
	grpcserver := grpc.NewServer()

	// register service
	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(grpcserver, &cs)

	// grpc lisnten and server
	err = grpcserver.Serve(listener)
	if err != nil {
		log.Fatalf("failed start grpc server: %v ", err)
	}

}
