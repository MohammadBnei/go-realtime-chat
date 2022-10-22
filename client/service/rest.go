package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"rc-client/client/chat"
	"rc-client/models"

	"astuart.co/go-sse"
)

type RestService interface {
	GetStream(host string, msg chan *sse.Event) error
	WriteData()
}

type restService struct {
	username string
	roomId   string
	host     string
	api      chat.ClientService
}

type RestServiceConfig struct {
	Username string
	RoomId   string
	Host     string
	Api      chat.ClientService
}

func NewRestService(config *RestServiceConfig) RestService {
	return &restService{
		username: config.Username,
		roomId:   config.RoomId,
		host:     config.Host,
		api:      config.Api,
	}
}

func (rs *restService) GetStream(host string, msg chan *sse.Event) error {
	err := sse.Notify(fmt.Sprintf("%v/api/stream/%v", rs.host, rs.roomId), msg)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (rs *restService) WriteData() {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error", err)
			break
		}
		if sendData == "" {
			break
		}
		sendData = sendData[:len(sendData)-1]
		postParams := chat.NewPostRoomIDParams()
		postParams.ID = rs.roomId
		postParams.MessageInput = &models.AdapterMessageInput{
			Data:   sendData,
			UserID: rs.username,
		}
		_, err = rs.api.PostRoomID(postParams)

		if err != nil {
			fmt.Println(err)
		}
	}
}
