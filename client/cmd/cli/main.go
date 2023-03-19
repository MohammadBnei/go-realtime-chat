package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"rc-client/domain"
	"rc-client/service"

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
	messages := make(chan *domain.Message, 100)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	api := messagegrpc.NewRoomClient(conn)
	chatService := service.NewGrpcService(api)

	username = StringPrompt("Username : ")
	roomId = StringPrompt("Room Id : ")

	go chatService.GetStream(roomId, messages)

	app := tview.NewApplication().EnableMouse(true)

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	list := tview.NewList().
		AddItem("Messages", "", rune(0), nil).
		SetMainTextColor(tcell.ColorGreenYellow).
		SetSecondaryTextColor(tcell.ColorDimGray)

	go func(list *tview.List) {
		for m := range messages {
			if m.UserId == username {
				list.AddItem(m.Text, fmt.Sprintf("<- %s", m.UserId), rune(0), nil)
			} else {
				list.AddItem(m.Text, fmt.Sprintf("-> %s", m.UserId), rune(0), nil)
			}
			app.Draw()
		}
	}(list)

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

	grid := tview.NewGrid().
		SetRows(2, 0, 2).
		SetBorders(true).
		AddItem(newPrimitive(fmt.Sprintf("Room\t: %s\nUsername\t: %s", roomId, username)), 0, 0, 1, 3, 0, 0, false).
		AddItem(inputField, 1, 0, 1, 1, 0, 0, false).
		AddItem(list, 1, 1, 1, 2, 0, 0, false)

	if err := app.SetRoot(grid, true).SetFocus(inputField).Run(); err != nil {
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
