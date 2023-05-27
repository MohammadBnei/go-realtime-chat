package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	adapter "github.com/MohammadBnei/go-realtime-chat/server/adapter/grpc"
	"github.com/MohammadBnei/go-realtime-chat/server/service"
	messagev1alpha "github.com/MohammadBnei/go-realtime-chat/server/stubs/message/v1alpha"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type config struct {
	secure bool
	port   int32
	cert   string
	key    string
}

func serveGrpc(conf *config) {
	roomManager := service.GetRoomManager()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", conf.port))
	if err != nil {
		log.Fatal(err)
	}

	var grpcServer *grpc.Server
	if conf.secure {
		creds, err := loadTLSCredentials(conf)
		if err != nil {
			log.Fatal(err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
			MinTime:             30 * time.Second,
		}))
	} else {
		grpcServer = grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
			MinTime:             30 * time.Second,
		}))
	}

	server := adapter.NewMessageAdapter(roomManager)

	messagev1alpha.RegisterMessageServiceServer(grpcServer, server)
	// roomgrpc.RegisterRoomServiceServer(grpcServer, server)
	reflection.Register(grpcServer)
	go func() {
		if conf.secure {
			log.Printf("Secure gRPC Server Started on port %v", conf.port)
		} else {
			log.Printf("Insecure gRPC Server Started on port %v", conf.port)
		}
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

func loadTLSCredentials(conf *config) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(conf.cert, conf.key)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
