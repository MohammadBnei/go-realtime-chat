package cmd

import (
	"crypto/tls"
	"log"
	"strings"

	window "github.com/MohammadBnei/go-realtime-chat/client/cli"
	"github.com/MohammadBnei/go-realtime-chat/client/domain"
	"github.com/MohammadBnei/go-realtime-chat/client/service"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"github.com/johnsiilver/getcert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func cli(conf *domain.Config) {
	var conn *grpc.ClientConn

	creds := insecure.NewCredentials()

	if conf.Secure {
		tlsCert, _, err := getcert.FromTLSServer(conf.Host, false)
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		servName := strings.Split(conf.Host, ":")[0]
		creds = credentials.NewTLS(&tls.Config{
			ServerName:   servName,
			Certificates: []tls.Certificate{tlsCert},
		})
	}
	conn, err := grpc.Dial(conf.Host, grpc.WithTransportCredentials(creds), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time: 30,
	}))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	messages := make(chan *domain.Message, 100)
	panicChan := make(chan error)
	quitChan := make(chan bool)

	api := messagegrpc.NewRoomClient(conn)
	chatService := service.NewGrpcService(api, panicChan, quitChan)

	getStream := func(messages chan *domain.Message) func(roomId string) {
		return func(roomId string) {
			go chatService.GetStream(roomId, messages)

			if err := <-panicChan; err != nil {
				panic(err)
			}
		}
	}(messages)

	getStream(conf.Room)

	window.DrawWindow(chatService, conf, getStream, messages, quitChan)

}
