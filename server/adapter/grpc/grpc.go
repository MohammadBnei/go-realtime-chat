package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/MohammadBnei/go-realtime-chat/server/service"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"buf.build/gen/go/bneiconseil/go-chat/protocolbuffers/go/message"
)

type grpcAdapter struct {
	messagegrpc.UnimplementedRoomServer
	roomManager service.Manager
}

func NewGrpcAdapter(rm service.Manager) messagegrpc.RoomServer {
	return &grpcAdapter{roomManager: rm}
}

func (ga *grpcAdapter) GetRoom(ctx context.Context, rq *message.RoomRequest) (*message.RoomResponse, error) {
	return nil, nil
}

func (ga *grpcAdapter) PostToRoom(ctx context.Context, msg *message.Message) (*message.RoomResponse, error) {
	ga.roomManager.Submit(msg.UserId, msg.RoomId, msg.Text)
	logRequest(msg.RoomId, "Post to Room")
	return &message.RoomResponse{
		Success: true,
	}, nil
}
func (ga *grpcAdapter) DeleteRoom(ctx context.Context, rq *message.RoomRequest) (*message.RoomResponse, error) {
	return nil, nil
}

func (ga *grpcAdapter) StreamRoom(rr *message.RoomRequest, srs messagegrpc.Room_StreamRoomServer) error {
	listener := ga.roomManager.OpenListener(rr.RoomId)
	defer ga.roomManager.CloseListener(rr.RoomId, listener)
	logRequest(rr.RoomId, "Stream Room")

	for {
		select {
		case msg := <-listener:
			serviceMsg, ok := msg.(service.Message)
			if !ok {
				fmt.Println(msg)
				continue
			}

			if err := srs.Send(&message.Message{
				UserId: serviceMsg.UserId,
				RoomId: serviceMsg.RoomId,
				Text:   serviceMsg.Text,
			}); err != nil {
				return err
			}
		case <-srs.Context().Done():
			return nil
		}
	}
}

func logRequest(room, method string) {
	time := time.Now().Format("2006-01-02 - 15:04:05")

	fmt.Printf("[GRPC] %s %s room : %s \n", time, method, room)
}
