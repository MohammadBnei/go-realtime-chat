package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"rc-client/domain"
	"rc-client/messagePB"
)

type Service interface {
	GetStream(host string, msg chan *domain.Message) error
	WriteData()
}

type grpcService struct {
	username, roomId, host string
	api                    messagePB.RoomClient
}
type GrpcServiceConfig struct {
	Username string
	RoomId   string
	Host     string
	Api      messagePB.RoomClient
}

func NewGrpcService(config *GrpcServiceConfig) Service {
	return &grpcService{config.Username, config.RoomId, config.Host, config.Api}
}

func (rs *grpcService) GetStream(host string, msg chan *domain.Message) error {
	src, err := rs.api.StreamRoom(context.Background(), &messagePB.RoomRequest{
		UserId: rs.username,
		RoomId: rs.roomId,
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

func (rs *grpcService) WriteData() {
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
		sendData = sendData[:len(sendData)-1]
		res, err := rs.api.PostToRoom(context.Background(), &messagePB.Message{
			UserId: rs.username,
			RoomId: rs.roomId,
			Text:   sendData,
		})
		if err != nil {
			fmt.Println(err)
		}

		if !res.Success {
			fmt.Println("Something went wrong")
		}
	}
}
