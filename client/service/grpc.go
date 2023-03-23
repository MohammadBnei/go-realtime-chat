package service

import (
	"context"
	"fmt"
	"io"

	"github.com/MohammadBnei/go-realtime-chat/client/domain"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"buf.build/gen/go/bneiconseil/go-chat/protocolbuffers/go/message"
)

type Service interface {
	GetStream(roomId string, msg chan *domain.Message) error
	PostMessage(username, roomId, text string)
}

type grpcService struct {
	api       messagegrpc.RoomClient
	panicChan chan error
	quitChan  chan bool
}

func NewGrpcService(api messagegrpc.RoomClient, panicChan chan error, quitChan chan bool) Service {
	return &grpcService{api, panicChan, quitChan}
}

func (rs *grpcService) GetStream(roomId string, msg chan *domain.Message) error {
	streamClient, err := rs.api.StreamRoom(context.Background(), &message.RoomRequest{
		RoomId: roomId,
	})
	if err != nil {
		rs.panicChan <- err
	}

	rs.panicChan <- nil

	serverMsgChannel := make(chan *message.Message)
	errorChannel := make(chan error)
	toChannel := func(streamClient messagegrpc.Room_StreamRoomClient, serverMsgChannel chan *message.Message, errorChannel chan error) {
		for {
			newMsg, err := streamClient.Recv()
			if err != nil {
				errorChannel <- err
				return
			}
			serverMsgChannel <- newMsg
		}
	}

	go toChannel(streamClient, serverMsgChannel, errorChannel)

	for {
		select {
		case <-rs.quitChan:
			return nil
		case err := <-errorChannel:
			if err == io.EOF {
				return streamClient.CloseSend()
			}
			if err != nil {
				fmt.Println(err)
				return err
			}
		case newMsg := <-serverMsgChannel:
			msg <- &domain.Message{
				UserId: newMsg.UserId,
				RoomId: newMsg.RoomId,
				Text:   newMsg.Text,
			}
		}
	}
}

func (rs *grpcService) PostMessage(username, roomId, text string) {
	rs.api.PostToRoom(context.Background(), &message.Message{
		UserId: username,
		RoomId: roomId,
		Text:   text,
	})
}
