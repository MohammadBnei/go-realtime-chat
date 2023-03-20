package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	adapter "realtime-chat/adapter/grpc"
	"realtime-chat/config"
	"realtime-chat/service"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	roomManager := service.GetRoomManager()

	lis, err := net.Listen("tcp", "0.0.0.0:"+config.Config.Port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	server := adapter.NewGrpcAdapter(roomManager)

	messagegrpc.RegisterRoomServer(grpcServer, server)
	reflection.Register(grpcServer)
	go func() {
		log.Println("gRPC Server Started on : " + config.Config.Port)
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
