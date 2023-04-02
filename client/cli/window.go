package cli

import (
	"fmt"

	"github.com/MohammadBnei/go-realtime-chat/client/domain"
	"github.com/MohammadBnei/go-realtime-chat/client/service"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DrawWindow(chatService service.Service, conf *domain.Config, getStream func(username, roomId string), messages chan *domain.Message, quitChan chan bool) {
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
				chatService.PostMessage(conf.Username, conf.Room, inputField.GetText())
				inputField.SetText("")
			}
		})

	changeRoom := func(form *tview.Form) {
		if newRoomId := form.GetFormItemByLabel("Room Id").(*tview.InputField).GetText(); newRoomId != conf.Room {
			conf.Room = form.GetFormItemByLabel("Room Id").(*tview.InputField).GetText()
			quitChan <- true
			getStream(conf.Username, conf.Room)
			app.SetFocus(inputField)
		}
	}
	changeUsername := func(form *tview.Form) {
		if newUsername := form.GetFormItemByLabel("Username").(*tview.InputField).GetText(); newUsername != conf.Username {
			conf.Username = form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			quitChan <- true
			getStream(conf.Username, conf.Room)
			app.SetFocus(inputField)
		}
	}

	form := tview.NewForm().
		AddInputField("Username", conf.Username, 20, nil, nil).
		AddInputField("Room Id", conf.Room, 30, nil, nil)

	form.GetFormItemByLabel("Username").
		SetFinishedFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				form.GetFormItemByLabel("Username").(*tview.InputField).SetText(conf.Username)
			case tcell.KeyEnter:
				changeUsername(form)
			}
		})

	form.AddButton("Change username", func() {
		changeUsername(form)
	})
	form.GetFormItemByLabel("Room Id").
		SetFinishedFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEscape:
				form.GetFormItemByLabel("Room Id").(*tview.InputField).SetText(conf.Room)
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
			if m.UserId == conf.Username {
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
