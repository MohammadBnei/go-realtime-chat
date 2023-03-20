package main

import (
	"fmt"
	adapter "realtime-chat/adapter/html"
	"realtime-chat/config"
	"realtime-chat/service"

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

	router.Run(fmt.Sprintf(":%v", config.ParseConfig().Port))
}
