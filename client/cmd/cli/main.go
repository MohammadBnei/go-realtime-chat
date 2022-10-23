package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"rc-client/domain"
	"rc-client/messagePB"
	"rc-client/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var username, roomId string

func main() {
	host := os.Getenv("CHAT_HOST")
	if host == "" {
		host = "localhost:4002"
	}
	message := make(chan *domain.Message, 100)

	for username == "" {
		fmt.Print("Enter your username : ")
		_, err := fmt.Scan(&username)
		if err != nil {
			fmt.Println(err)
		}
	}
	for roomId == "" {

		fmt.Print("Enter the room name : ")
		_, err := fmt.Scan(&roomId)
		if err != nil {
			fmt.Println(err)
		}
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	api := messagePB.NewRoomClient(conn)
	restService := service.NewGrpcService(&service.GrpcServiceConfig{
		Username: username,
		RoomId:   roomId,
		Host:     "http://" + host,
		Api:      api,
	})

	go restService.GetStream("http://"+host, message)
	log.Println("Listening for stream")

	go restService.WriteData(os.Stdin)

	go printMessages(message)

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	// Stop the server
	log.Println("\nStopping stream")
}

func printMessages(message chan *domain.Message) {
	for {
		msg, ok := <-message
		if !ok {
			fmt.Println("Channel closed")
			break
		}
		if msg.UserId == username {
			continue
		}
		if string(msg.Text) != "\n" {
			// if msg.UserId != username {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m → %s\n> ", msg.UserId, msg.Text)
			// }
		}
	}
}
