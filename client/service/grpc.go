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
	GetStream(username, roomId string, msg chan *domain.Message) error
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

func (rs *grpcService) GetStream(username, roomId string, msg chan *domain.Message) error {
	ctx, cancel := context.WithCancel(context.Background())
	streamClient, err := rs.api.StreamRoom(ctx, &message.RoomRequest{
		RoomId: roomId,
		UserId: username,
	})
	defer cancel()
	if err != nil {
		rs.panicChan <- err
		return err
	}

	rs.panicChan <- nil

	serverMsgChannel := make(chan *message.Message)
	errorChannel := make(chan error)

	go func(streamClient messagegrpc.Room_StreamRoomClient, serverMsgChannel chan *message.Message, errorChannel chan error) {
		for {
			newMsg, err := streamClient.Recv()
			if err != nil {
				errorChannel <- err
				return
			}
			serverMsgChannel <- newMsg
		}
	}(streamClient, serverMsgChannel, errorChannel)

	for {
		select {
		case <-rs.quitChan:
			cancel()
			return nil
		case err := <-errorChannel:
			if err == io.EOF {
				return err
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
