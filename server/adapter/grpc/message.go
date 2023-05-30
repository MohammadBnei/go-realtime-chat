package adapter

import (
	"fmt"
	"time"

	"github.com/MohammadBnei/go-realtime-chat/server/service"
	messagev1alpha "github.com/MohammadBnei/go-realtime-chat/server/stubs/message/v1alpha"
)

type messageAdapter struct {
	messagev1alpha.UnimplementedMessageServiceServer
	roomManager service.Manager
}

func NewMessageAdapter(rm service.Manager) messagev1alpha.MessageServiceServer {
	return &messageAdapter{roomManager: rm}
}

func (a *messageAdapter) StreamMessages(streamServer messagev1alpha.MessageService_StreamMessagesServer) error {
	sendChannel := make(chan string, 100)
	configChannel := make(chan *messagev1alpha.Config, 1)

	go func() {
		for {
			msg, err := streamServer.Recv()
			if err != nil {
				fmt.Println(err)
				return
			}
			// TODO: Check passphrase

			switch reflectMsg := msg.Message.(type) {
			case *messagev1alpha.StreamMessagesRequest_Text:
				msg.Config.Passphrase = ""
				configChannel <- msg.Config
				a.roomManager.Room(msg.Config.GetRoomId()).Submit(&messagev1alpha.StreamMessagesRequest{
					Message: reflectMsg,
					Config:  msg.Config,
				})

			}

			fmt.Println(msg)
		}
	}()

	go func() {
		for msg := range sendChannel {
			streamServer.Send(
				&messagev1alpha.StreamMessagesResponse{
					Message: &messagev1alpha.StreamMessagesResponse_Text{Text: &messagev1alpha.TextMessage{Text: msg}},
				},
			)
		}
	}()
	defer close(sendChannel)

	return nil
}

func logRequest(room, method string) {
	time := time.Now().Format("2006-01-02 - 15:04:05")

	fmt.Printf("[GRPC] %s %s room : %s \n", time, method, room)
}
