package main

import (
	"fmt"

	adapter "github.com/MohammadBnei/realtime-chat/server/adapter/html"
	"github.com/MohammadBnei/realtime-chat/server/config"
	"github.com/MohammadBnei/realtime-chat/server/service"

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
