package adapter

import "realtime-chat/service"

type grpcAdapter struct {
	roomManager service.Manager
}

func NewGrpcAdapter(rm service.Manager) *grpcAdapter {
	return &grpcAdapter{rm}
}
