package cmd

import (
	"crypto/tls"
	"log"
	"strings"
	"time"

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
	ka := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time: 45 * time.Second,
	})

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
		ka = grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 45 * time.Second,
		})
	}
	conn, err := grpc.Dial(conf.Host, grpc.WithTransportCredentials(creds), ka)
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	messages := make(chan *domain.Message, 100)
	panicChan := make(chan error)
	quitChan := make(chan bool)

	api := messagegrpc.NewRoomClient(conn)
	chatService := service.NewGrpcService(api, panicChan, quitChan)

	getStream := func(messages chan *domain.Message, quitChan chan bool) func(username, roomId string) {
		return func(username, roomId string) {
			go chatService.GetStream(username, roomId, messages)

			if err := <-panicChan; err != nil {
				panic(err)
			}
		}
	}(messages, quitChan)

	getStream(conf.Username, conf.Room)

	window.DrawWindow(chatService, conf, getStream, messages, quitChan)

}
