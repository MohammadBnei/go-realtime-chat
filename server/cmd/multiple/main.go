package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	html "realtime-chat/adapter/html"
	rest "realtime-chat/adapter/rest"
	"realtime-chat/cmd/rest/docs"
	"realtime-chat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	go StartRest("4001")
	go StartHTML("4000")

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	// Stop the server
	log.Println("server stopped")
}

func StartHTML(port string) {
	roomManager := service.GetRoomManager()
	adapter := html.NewGinHTMLAdapter(roomManager)
	router := gin.Default()
	router.SetHTMLTemplate(adapter.Template)

	router.GET("/room/:roomid", adapter.GetRoom)
	router.POST("/room/:roomid", adapter.PostRoom)
	router.DELETE("/room/:roomid", adapter.DeleteRoom)
	router.GET("/stream/:roomid", adapter.Stream)

	router.Run(fmt.Sprintf(":%v", port))
}

func StartRest(port string) {
	roomManager := service.GetRoomManager()
	adapter := rest.NewGinAdapter(roomManager)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"

	api := router.Group("/api")
	api.POST("/room/:roomid", adapter.PostRoom)
	api.DELETE("/room/:roomid", adapter.DeleteRoom)
	api.GET("/stream/:roomid", adapter.Stream)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(fmt.Sprintf(":%v", port))
}
