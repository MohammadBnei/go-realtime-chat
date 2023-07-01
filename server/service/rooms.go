package service

import "github.com/MohammadBnei/go-realtime-chat/server/broadcast"

type Manager[T any] interface {
	OpenListener(string) chan interface{}
	CloseListener(string, chan interface{})
	Submit(string, string, T)
	DeleteBroadcast(string)
}

type Message[T any] struct {
	UserId  string
	RoomId  string
	Content T
}

type Listener struct {
	RoomId string
	Chan   chan interface{}
}

type manager[T any] struct {
	roomChannels map[string]broadcast.Broadcaster[T]
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message[T]
}

func NewRoomManager[T any]() Manager[T] {
	managerSingleton := &manager[T]{
		roomChannels: make(map[string]broadcast.Broadcaster[T]),
		open:         make(chan *Listener, 100),
		close:        make(chan *Listener, 100),
		delete:       make(chan string, 100),
		messages:     make(chan *Message[T], 100),
	}

	go managerSingleton.run()

	return managerSingleton
}

func (m *manager[T]) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case roomid := <-m.delete:
			m.deleteBroadcast(roomid)
		case msg := <-m.messages:
			m.room(msg.RoomId).Submit(msg.Content)
		}
	}
}

func (m *manager[T]) register(listener *Listener) {
	m.room(listener.RoomId).Register(listener.Chan)
}

func (m *manager[T]) deregister(listener *Listener) {
	m.room(listener.RoomId).Unregister(listener.Chan)
	close(listener.Chan)
}

func (m *manager[T]) deleteBroadcast(roomid string) {
	b, ok := m.roomChannels[roomid]
	if ok {
		b.Close()
		delete(m.roomChannels, roomid)
	}
}

/*
Get the room with the id roomid, or creates and registers it
*/
func (m *manager[T]) room(roomid string) broadcast.Broadcaster[T] {
	b, ok := m.roomChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster[T](10)
		m.roomChannels[roomid] = b
	}
	return b
}

func (m *manager[T]) OpenListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &Listener{
		RoomId: roomid,
		Chan:   listener,
	}
	return listener
}

func (m *manager[T]) CloseListener(roomid string, channel chan interface{}) {
	m.close <- &Listener{
		RoomId: roomid,
		Chan:   channel,
	}
}

func (m *manager[T]) DeleteBroadcast(roomid string) {
	m.delete <- roomid
}

func (m *manager[T]) Submit(userId, roomId string, msg T) {
	m.messages <- &Message[T]{
		UserId:  userId,
		RoomId:  roomId,
		Content: msg,
	}
}
