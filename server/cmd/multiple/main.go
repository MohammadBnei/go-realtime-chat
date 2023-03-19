package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"buf.build/gen/go/bneiconseil/go-chat/grpc/go/message/messagegrpc"
	adapter "github.com/MohammadBnei/realtime-chat/server/adapter/grpc"
	html "github.com/MohammadBnei/realtime-chat/server/adapter/html"
	rest "github.com/MohammadBnei/realtime-chat/server/adapter/rest"
	"github.com/MohammadBnei/realtime-chat/server/cmd/rest/docs"
	"github.com/MohammadBnei/realtime-chat/server/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	go StartRest("4001")
	go StartHTML("4000")
	go StartGrpc("4002")

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

func StartGrpc(port string) {
	roomManager := service.GetRoomManager()

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	server := adapter.NewGrpcAdapter(roomManager)

	messagegrpc.RegisterRoomServer(grpcServer, server)
	reflection.Register(grpcServer)
	go func() {
		log.Println("gRPC Server Started on : " + port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
