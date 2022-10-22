package main

import (
	"fmt"
	"log"
	"os"

	"rc-client/client/chat"
	"rc-client/service"

	sse "astuart.co/go-sse"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

var username, roomId string

func main() {
	host := os.Getenv("CHAT_HOST")
	if host == "" {
		host = "localhost:4001"
	}
	message := make(chan *sse.Event, 100)

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

	api := chat.New(httptransport.New(host, "/api", nil), strfmt.Default)
	restService := service.NewRestService(&service.RestServiceConfig{
		Username: username,
		RoomId:   roomId,
		Host:     "http://" + host,
		Api:      api,
	})

	go restService.GetStream("http://"+host, message)
	log.Println("Listening for stream")

	go restService.WriteData()

	go func() {
		for {
			msg, ok := <-message
			if !ok {
				fmt.Println("Channel closed")
				break
			}
			str := make([]byte, 1024)
			n, err := msg.Data.Read(str)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				continue
			}
			if string(str) != "\n" {
				// Green console colour: 	\x1b[32m
				// Reset console colour: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\n\x1b[0m> ", str)
			}
		}
	}()

	select {}
}
