package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"rc-client/domain"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"buf.build/gen/go/bneiconseil/go-chat/protocolbuffers/go/message"
)

type Service interface {
	GetStream(oomId string, msg chan *domain.Message) error
	WriteData(username, roomId string, rw io.Reader)
}

type grpcService struct {
	host string
	api  messagegrpc.RoomClient
}
type GrpcServiceConfig struct {
	Host string
	Api  messagegrpc.RoomClient
}

func NewGrpcService(config *GrpcServiceConfig) Service {
	return &grpcService{config.Host, config.Api}
}

func (rs *grpcService) GetStream(roomId string, msg chan *domain.Message) error {
	src, err := rs.api.StreamRoom(context.Background(), &message.RoomRequest{
		RoomId: roomId,
	})
	if err != nil {
		return err
	}

	for {
		newMsg, err := src.Recv()
		if err == io.EOF {
			return src.CloseSend()
		}
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
		res, err := rs.api.PostToRoom(context.Background(), &message.Message{
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
