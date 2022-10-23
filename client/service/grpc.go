package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"rc-client/domain"
	"rc-client/messagePB"
)

type Service interface {
	GetStream(oomId string, msg chan *domain.Message) error
	WriteData(username, roomId string, rw io.Reader)
}

type grpcService struct {
	host string
	api  messagePB.RoomClient
}
type GrpcServiceConfig struct {
	Host string
	Api  messagePB.RoomClient
}

func NewGrpcService(config *GrpcServiceConfig) Service {
	return &grpcService{config.Host, config.Api}
}

func (rs *grpcService) GetStream(roomId string, msg chan *domain.Message) error {
	src, err := rs.api.StreamRoom(context.Background(), &messagePB.RoomRequest{
		RoomId: roomId,
	})
	if err != nil {
		return err
	}

	for {
		var newMsg messagePB.Message
		err := src.RecvMsg(&newMsg)
		if err != nil {
			fmt.Println(err)
			return err
		}

		msg <- &domain.Message{
			UserId: newMsg.UserId,
			RoomId: newMsg.RoomId,
			Text:   newMsg.Text,
		}
	}
}

func (rs *grpcService) WriteData(username, roomId string, rw io.Reader) {
	stdReader := bufio.NewReader(rw)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error", err)
			break
		}
		if sendData == "" || sendData == "\n" {
			continue
		}
		sendData = sendData[:len(sendData)-1]
		res, err := rs.api.PostToRoom(context.Background(), &messagePB.Message{
			UserId: username,
			RoomId: roomId,
			Text:   sendData,
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !res.Success {
			fmt.Println("Something went wrong")
		}
	}
}
