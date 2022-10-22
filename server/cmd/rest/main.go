package main

import (
	"fmt"
	adapter "realtime-chat/adapter/rest"
	"realtime-chat/cmd/rest/docs"
	"realtime-chat/config"
	"realtime-chat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Realtime Chat API
// @version         0.1
// @description     Realtime chat api using channels.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api

func main() {
	roomManager := service.GetRoomManager()
	adapter := adapter.NewGinAdapter(roomManager)
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"

	api := router.Group("/api")
	api.POST("/room/:roomid", adapter.PostRoom)
	api.DELETE("/room/:roomid", adapter.DeleteRoom)
	api.GET("/stream/:roomid", adapter.Stream)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(fmt.Sprintf(":%v", config.ParseConfig().ServerConfig.Port))
}
