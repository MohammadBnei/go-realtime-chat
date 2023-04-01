package cmd

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MohammadBnei/go-realtime-chat/client/domain"
	"github.com/MohammadBnei/go-realtime-chat/client/service"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	"github.com/gdamore/tcell/v2"
	"github.com/johnsiilver/getcert"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type config struct {
	host     string
	secure   bool
	username string
	room     string
}

func cli(conf *config) {
	var conn *grpc.ClientConn

	creds := insecure.NewCredentials()

	if conf.secure {
		tlsCert, _, err := getcert.FromTLSServer(conf.host, false)
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		servName := strings.Split(conf.host, ":")[0]
		creds = credentials.NewTLS(&tls.Config{
			ServerName:   servName,
			Certificates: []tls.Certificate{tlsCert},
		})
	}
	conn, err := grpc.Dial(conf.host, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	messages := make(chan *domain.Message, 100)
	panicChan := make(chan error)
	quitChan := make(chan bool)

	api := messagegrpc.NewRoomClient(conn)
	chatService := service.NewGrpcService(api, panicChan, quitChan)

	if conf.username == "" {
		conf.username = StringPrompt("Username : ")
	}

	getStream := func(messages chan *domain.Message) func(roomId string) {
		return func(roomId string) {
			go chatService.GetStream(roomId, messages)

			if err := <-panicChan; err != nil {
				panic(err)
			}
		}
	}(messages)

	getStream(conf.room)

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
				chatService.PostMessage(conf.username, conf.room, inputField.GetText())
				inputField.SetText("")
			}
		})

	changeRoom := func(form *tview.Form) {
		if newRoomId := form.GetFormItemByLabel("Room Id").(*tview.InputField).GetText(); newRoomId != conf.room {
			quitChan <- true
			go chatService.PostMessage(conf.username, conf.room, fmt.Sprintf("%s disconnected", conf.username))
			conf.room = form.GetFormItemByLabel("Room Id").(*tview.InputField).GetText()
			getStream(conf.room)
			app.SetFocus(inputField)
			go chatService.PostMessage(conf.username, conf.room, fmt.Sprintf("*%s connected*", conf.username))
		}
	}

	form := tview.NewForm().
		AddInputField("Username", conf.username, 20, nil, func(text string) { conf.username = text }).
		AddInputField("Room Id", conf.room, 30, nil, nil)

	form.GetFormItemByLabel("Room Id").
		SetFinishedFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				form.GetFormItemByLabel("Room Id").(*tview.InputField).SetText(conf.room)
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
			if m.UserId == conf.username {
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

	go chatService.PostMessage(conf.username, conf.room, fmt.Sprintf("*%s connected to %s*", conf.username, conf.room))

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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("localhost.pem", "localhost-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
