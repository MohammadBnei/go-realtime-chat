package adapter

import (
	"context"

	"github.com/MohammadBnei/go-realtime-chat/server/service"

	roomv1alpha "github.com/MohammadBnei/go-realtime-chat/server/stubs/room/v1alpha"
)

type roomAdapter struct {
	roomv1alpha.UnimplementedRoomServiceServer
	roomManager service.Manager
}

func NewRoomAdapter(rm service.Manager) roomv1alpha.RoomServiceServer {
	return &roomAdapter{roomManager: rm}
}

func (s *roomAdapter) ListRooms(context.Context, *roomv1alpha.ListRoomsRequest) (*roomv1alpha.ListRoomsResponse, error) {
	return nil, nil
}

func (s *roomAdapter) CreateRoom(context.Context, *roomv1alpha.CreateRoomRequest) (*roomv1alpha.CreateRoomResponse, error) {
	return nil, nil
}

func (s *roomAdapter) ChangePassphrase(context.Context, *roomv1alpha.ChangePassphraseRequest) (*roomv1alpha.ChangePassphraseResponse, error) {
	return nil, nil
}

func (s *roomAdapter) UpdateRoomId(context.Context, *roomv1alpha.UpdateRoomIdRequest) (*roomv1alpha.UpdateRoomIdResponse, error) {
	return nil, nil
}

func (s *roomAdapter) DeleteRoom(context.Context, *roomv1alpha.DeleteRoomRequest) (*roomv1alpha.DeleteRoomResponse, error) {
	return nil, nil
}
