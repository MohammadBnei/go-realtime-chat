# Realtime Chat Server Tutorial

## Objectives

We want to create a realtime chat application. To do so, we will start with the server and implement the following functionalities :
- Broadcast system
  - take a message and broadcast it to all connected listeners
- Room system
  - One broadcaster per room, to have different group discussion possibility
- Adapter
  - How the services are exposed and initialized
  - gRPC
  - REST
  - HTML template
- Commad
  - How the application will be started
  - All adapters
  - Only one adapter


## Setup

To start, create a new directory and initialize a go module with your github path :
```bash
# Replace with your username
go mod init github.com/$GITHUB_USERNAME/realtime-chat-server
```

That's it, we are ready to start implementing

## Broadcaster

### Interface

The broadcaster interface tells us what it needs to do :

```golang
// The Broadcaster interface describes the main entry points to
// broadcasters.
type Broadcaster interface {

	// Register a new channel to receive broadcasts
	Register(chan<- interface{})

	// Unregister a channel so that it no longer receives broadcasts.
	Unregister(chan<- interface{})

	// Shut this broadcaster down.
	Close() error
	// Submit a new object to all subscribers
	Submit(interface{})
}
```

We will be using channels to coordinate the messages with the listeners, also to handle subscribtion and unregister.

So, let's create a broadcast/broadcaster.go file. Set the package at the top (*broadcast*) and write the Broadcaster interface.

### Struct

The struct will carry all the variables needed to make the broadcaster work. As mentionned, there is a lot of channels.
```golang
type broadcaster struct {
	input chan interface{}
	reg   chan chan<- interface{}
	unreg chan chan<- interface{}

	outputs map[chan<- interface{}]bool
}
```

The important thing to note here is the outputs : it's a map with channels as key and bool as value.

Next, let's implement the Factory. This is the initialization of our struct :
```golang
// NewBroadcaster creates a new broadcaster with the given input
// channel buffer length.
func NewBroadcaster(buflen int) Broadcaster {
	b := &broadcaster{
		input:   make(chan interface{}, buflen),
		reg:     make(chan chan<- interface{}),
		unreg:   make(chan chan<- interface{}),
		outputs: make(map[chan<- interface{}]bool),
	}

	go b.run()

	return b
}
```

This creates the necessary channels and pointers. Notice the extensive use of *make*, to allocate memory and other cool go stuff.

So, *reg* and *unreg* are channels of channel. They pass channel around. We handle the registration this way because the listeners are actually channels themself, so we cut unnecessary intermediate to use the listening channel directly.

Finally, we run the broadcaster in its own goroutine. This will free the main thread to do other things, and handle the message broadcasting in an efficient go way.

### Functions

The actual broadcast function will iterate over all the channels in the outputs map and forward the message to them : 
```golang
func (b *broadcaster) broadcast(m interface{}) {
	for ch := range b.outputs {
		ch <- m
	}
}
```

The central element to the broadcaster is the run function. It ties together all the channels and their goal, let's check it out :

```golang
func (b *broadcaster) run() {
	for {
		select {
    // On any input, broadcast to all registered listeners, aka outputs
		case m := <-b.input:
			b.broadcast(m)

    // Handle registration/unregistration of listeners by channels. 
		case ch, ok := <-b.reg:
			if ok {
				b.outputs[ch] = true
			} else {
				return
			}
		case ch := <-b.unreg:
			delete(b.outputs, ch)
		}
	}
}
```

Remember we said *reg* is a channel of channel ? So, when a new channel is registered, this select will make sure there will be no race conditions for the broadcast.

Now, some easy stuff. To register or register a listener, we only now have to pass the listener (in the form of a channel) to the *reg*/*unreg* channel :

```golang
func (b *broadcaster) Register(newch chan<- interface{}) {
	b.reg <- newch
}

func (b *broadcaster) Unregister(newch chan<- interface{}) {
	b.unreg <- newch
}
```

The close function stops the broadcaster from adding or removing listeners :
```golang
func (b *broadcaster) Close() error {
	close(b.reg)
	close(b.unreg)
	return nil
}
```

And lastly, the submit :
```golang
// Submit an item to be broadcast to all listeners.
func (b *broadcaster) Submit(m interface{}) {
	if b != nil {
		b.input <- m
	}
}
```

That's it ! We have our broadcasting block. Next, we will create the room service that will handle broadcasting per room.

## Room Service

As usual, let's start with the interface to grasp the wanted behavior. Create the service/rooms.go file and set the correct package. Then, write the following :
```golang
type Manager interface {
	OpenListener(roomid string) chan interface{}
	CloseListener(roomid string, channel chan interface{})
	Submit(userid, roomid, text string)
	DeleteBroadcast(roomid string)
}
```

Our room manager is simple. On *open/closeListener*, we want to add/remove a listening channel from the room with *roomid*. We will here create the room if needed.

On submit, we will broadcast the *text* message and the *userid* on all listeners of the room with id *roomid*. We will be using an internal channel for this, *messages*. 

DeleteBroadcast is self explanatory.

We will need some structs to hold our informations and internal variables :

```golang
type Message struct {
	UserId string
	RoomId string
	Text   string
}

type Listener struct {
	RoomId string
	Chan   chan interface{}
}

type manager struct {
	roomChannels map[string]broadcast.Broadcaster
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message
}
```

The *Message* and *Listener* are dictionnary like entities, and the *manage* is our class. We choose not to strongly tie the listener to the room, only by the *roomid* string. 

Here is *Open* and *CloseListener* :
```golang
func (m *manager) OpenListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &Listener{
		RoomId: roomid,
		Chan:   listener,
	}
	return listener
}

func (m *manager) CloseListener(roomid string, channel chan interface{}) {
	m.close <- &Listener{
		RoomId: roomid,
		Chan:   channel,
	}
}
```

Like broadcaster, we will use channel to handle the internal communication of the service. Because of that, our public facing code is easy implemented because the complexity is transfered to internal service logic.

Here is the delete broadcast and submit functions :

```golang
func (m *manager) DeleteBroadcast(roomid string) {
	m.delete <- roomid
}

func (m *manager) Submit(userid, roomid, text string) {
	msg := &Message{
		UserId: userid,
		RoomId: roomid,
		Text:   text,
	}
	m.messages <- msg
}
```

Now that we have every subpart, let's get into the main logic :
```golang
func (m *manager) register(listener *Listener) {
	m.room(listener.RoomId).Register(listener.Chan)
}

func (m *manager) deregister(listener *Listener) {
	m.room(listener.RoomId).Unregister(listener.Chan)
	close(listener.Chan)
}

func (m *manager) deleteBroadcast(roomid string) {
	b, ok := m.roomChannels[roomid]
	if ok {
		b.Close()
		delete(m.roomChannels, roomid)
	}
}

/*
Get the room with the id roomid, or creates and registers it
*/
func (m *manager) room(roomid string) broadcast.Broadcaster {
	b, ok := m.roomChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.roomChannels[roomid] = b
	}
	return b
}
```
These functions handle the link between our service and the broadcaster. They are straighforward.

The *run* function will tie our channels together :
```golang
func (m *manager) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case roomid := <-m.delete:
			m.deleteBroadcast(roomid)
		case message := <-m.messages:
			m.room(message.RoomId).Submit(*message)
		}
	}
}
```

Let's think about the usefulness of channels here. We have a system that will eventually handle thousands of messages per seconds, while allowing clients to open rooms, register or unregister.

But we have to be carefull. If a message is sent to a closed channel, golang will panic. So, we designed a system that keeps in sync by having one unique entry point, our room manager. 

```golang
func GetRoomManager() Manager {
	if managerSingleton == nil {
		managerSingleton = &manager{
			roomChannels: make(map[string]broadcast.Broadcaster),
			open:         make(chan *Listener, 100),
			close:        make(chan *Listener, 100),
			delete:       make(chan string, 100),
			messages:     make(chan *Message, 100),
		}

		go managerSingleton.run()
	}

	return managerSingleton
}
```

Lastly, we will make a singleton out of the roomManager. This means that every instance that calls the service will receive the same pointer. Let's look at it :
```golang
var managerSingleton *manager

func GetRoomManager() Manager {
	if managerSingleton == nil {
		managerSingleton = &manager{
			roomChannels: make(map[string]broadcast.Broadcaster),
			open:         make(chan *Listener, 100),
			close:        make(chan *Listener, 100),
			delete:       make(chan string, 100),
			messages:     make(chan *Message, 100),
		}

		go managerSingleton.run()
	}

	return managerSingleton
}
```

To test it out, we've provided you with a basic gin HTML package. Create a main.go at the root and write the following :
```golang
package main

import (
	"fmt"

	adapter "github.com/MohammadBnei/go-html-adapter/adapterHTML"
	"$YOUR_GOMODULE/service"

	"github.com/gin-gonic/gin"
)

var roomManager service.Manager

func main() {
	roomManager = service.GetRoomManager()
	adapter := adapter.NewGinHTMLAdapter(roomManager)
	router := gin.Default()
	router.SetHTMLTemplate(adapter.Template)

	router.GET("/room/:roomid", adapter.GetRoom)
	router.POST("/room/:roomid", adapter.PostRoom)
	router.DELETE("/room/:roomid", adapter.DeleteRoom)
	router.GET("/stream/:roomid", adapter.Stream)

	router.Run(fmt.Sprintf(":%v", 8080))
}
``` 
Start it and open your browser to [http://localhost:8080/room/test]. You can open a second tab and verify that the messages are correctly sent and received.

It's done ! Now you do the rest.

## REST Adapter

Implement a REST adapter and the main.go entrypoint.
We advise you to use [go-gin](https://gin-gonic.com), it will be easy to code SSE for client side streaming.

You will find the HTML template [here](/server/adapter/html/template.go). It will give you a client to quickly test your application.

You have to create 2 routes :
 - The **stream** route, which will take a *roomId* and stream messages
 - the **submit** route, which will post messages. Each message will have a *userId*, *roomId*, and *text*.