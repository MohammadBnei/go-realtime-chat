package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MohammadBnei/go-realtime-chat/client/domain"
	"github.com/MohammadBnei/go-realtime-chat/client/service"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"github.com/gdamore/tcell/v2"
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
	roomId = "lobby"
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	messages := make(chan *domain.Message, 100)
	panicChan := make(chan error)
	quitChan := make(chan bool)

	api := messagegrpc.NewRoomClient(conn)
	chatService := service.NewGrpcService(api, panicChan, quitChan)

	username = StringPrompt("Username : ")

	getStream := func(messages chan *domain.Message) func(roomId string) {
		return func(roomId string) {
			go chatService.GetStream(roomId, messages)

			if err := <-panicChan; err != nil {
				panic(err)
			}
		}
	}(messages)

	getStream(roomId)

	app := tview.NewApplication().EnableMouse(true)

	inputField := tview.NewInputField().
		SetLabel("Enter a message: ").
		SetFieldWidth(255)

	inputField.
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				app.Stop()
			case tcell.KeyEnter:
				chatService.PostMessage(username, roomId, inputField.GetText())
				inputField.SetText("")
			}
		})

	changeRoom := func(form *tview.Form) {
		if newRoomId := form.GetFormItemByLabel("Room Id").(*tview.InputField).GetText(); newRoomId != roomId {
			quitChan <- true
			chatService.PostMessage(username, roomId, fmt.Sprintf("%s disconnected", username))
			roomId = form.GetFormItemByLabel("Room Id").(*tview.InputField).GetText()
			getStream(roomId)
			app.SetFocus(inputField)
		}
	}

	form := tview.NewForm().
		AddInputField("Username", username, 20, nil, func(text string) { username = text }).
		AddInputField("Room Id", roomId, 30, nil, nil)

	form.GetFormItemByLabel("Room Id").
		SetFinishedFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				form.GetFormItemByLabel("Room Id").(*tview.InputField).SetText(roomId)
			case tcell.KeyEnter:
				changeRoom(form)
			}
		})

	form.AddButton("Change Room", func() {
		changeRoom(form)
	})

	list := tview.NewList().
		AddItem("Messages", "", rune(0), nil).
		SetMainTextColor(tcell.ColorGreenYellow).
		SetSecondaryTextColor(tcell.ColorBlanchedAlmond).
		SetWrapAround(false)

	list.
		Box.
		SetTitle("Messages").
		SetBackgroundColor(tcell.ColorDimGray)

	go func(list *tview.List) {
		for m := range messages {
			if m.UserId == username {
				list.AddItem(m.Text, fmt.Sprintf("-> %s", m.UserId), rune(0), nil).SetCurrentItem(-1)
			} else {
				list.AddItem(m.Text, fmt.Sprintf("<- %s", m.UserId), rune(0), nil).SetCurrentItem(-1)
			}
			app.Draw()
		}
	}(list)

	grid := tview.NewGrid().
		SetRows(0, 2).
		SetBorders(true).
		AddItem(form, 0, 0, 1, 1, 0, 0, false).
		AddItem(list, 0, 1, 1, 2, 0, 0, false).
		AddItem(inputField, 1, 0, 1, 3, 0, 0, true)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

}

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
