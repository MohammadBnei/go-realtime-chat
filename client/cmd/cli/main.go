package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"rc-client/client/chat"
	"rc-client/models"

	sse "astuart.co/go-sse"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

func main() {
	host := os.Getenv("CHAT_HOST")
	if host == "" {
		host = "localhost:4001"
	}
	message := make(chan *sse.Event, 100)

	go getStream("http://"+host, "test", message)
	log.Println("Listening for stream")

	api := chat.New(httptransport.New(host, "/api", nil), strfmt.Default)
	go writeData(api, "client-test", "test")

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

func getStream(host, roomId string, msg chan *sse.Event) error {
	err := sse.Notify(fmt.Sprintf("%v/api/stream/%v", host, roomId), msg)
	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func writeData(api chat.ClientService, username, roomId string) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error", err)
			break
		}
		if sendData == "" {
			break
		}
		postParams := chat.NewPostRoomIDParams()
		postParams.ID = roomId
		postParams.MessageInput = &models.AdapterMessageInput{
			Data:   sendData,
			UserID: username,
		}
		_, err = api.PostRoomID(postParams)

		if err != nil {
			fmt.Println(err)
		}
	}
}
