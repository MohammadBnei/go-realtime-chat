package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	adapter "github.com/MohammadBnei/realtime-chat/server/adapter/grpc"
	"github.com/MohammadBnei/realtime-chat/server/config"
	"github.com/MohammadBnei/realtime-chat/server/service"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	roomManager := service.GetRoomManager()
	config := config.ParseConfig()

	lis, err := net.Listen("tcp", "0.0.0.0:"+config.Port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	server := adapter.NewGrpcAdapter(roomManager)

	messagegrpc.RegisterRoomServer(grpcServer, server)
	reflection.Register(grpcServer)
	go func() {
		log.Println("gRPC Server Started on : " + config.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	// Stop the server
	log.Println("stopping the server")
	grpcServer.Stop()
	log.Println("server stopped")
}
