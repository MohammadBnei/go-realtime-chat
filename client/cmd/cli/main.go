package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"rc-client/domain"
	"rc-client/service"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"github.com/gosuri/uilive"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var username, roomId string

func main() {
	host := os.Getenv("CHAT_HOST")
	if host == "" {
		host = "localhost:4000"
	}
	message := make(chan *domain.Message, 100)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	api := messagegrpc.NewRoomClient(conn)
	grpcService := service.NewGrpcService(&service.GrpcServiceConfig{
		Host: "http://" + host,
		Api:  api,
	})

	go grpcService.WriteData("", "", os.Stdin)

	// go printMessages(message)

	go func() {
		newPrimitive := func(text string) tview.Primitive {
			return tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText(text)
		}
		main := newPrimitive("Main content")
		form := tview.NewForm().
			AddInputField("Username", "", 20, nil, func(text string) {
				username = text
			}).
			AddInputField("Room", "", 20, nil, func(text string) {
				roomId = text
			}).
			AddButton("Save", func() {
				go grpcService.GetStream(roomId, message)
			})

		grid := tview.NewGrid().
			SetRows(3, 0, 3).
			SetColumns(30, 0, 30).
			SetBorders(true).
			AddItem(newPrimitive(fmt.Sprintf("Room\t: %s\nUsername\t: %s", roomId, username)), 0, 0, 1, 3, 0, 0, false).
			AddItem(newPrimitive("Made with ❤️"), 2, 0, 1, 3, 0, 0, false)

		// Layout for screens narrower than 100 cells (menu and side bar are hidden).
		grid.AddItem(main, 1, 0, 1, 3, 0, 0, false)

		// Layout for screens wider than 100 cells.
		grid.AddItem(main, 1, 1, 1, 1, 0, 100, false).
			AddItem(form, 1, 2, 1, 1, 0, 100, false)

		if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	// Stop the server
	log.Println("\nStopping stream")
}

func printMessages(message chan *domain.Message) {
	writer := uilive.New()
	// start listening for updates and render
	writer.Start()
	for {
		msg, ok := <-message
		if !ok {
			fmt.Println("Channel closed")
			break
		}
		if msg.UserId == username {
			continue
		}
		if string(msg.Text) != "\n" {
			// if msg.UserId != username {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Fprintf(writer, "\x1b[32m%s\x1b[0m → %s\n> ", msg.UserId, msg.Text)
			// }
		}
	}
}
