package adapter

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"realtime-chat/service"

	"github.com/gin-gonic/gin"
)

type GinHTMLAdapter interface {
	Stream(c *gin.Context)
	GetRoom(c *gin.Context)
	PostRoom(c *gin.Context)
	DeleteRoom(c *gin.Context)
}

type ginHTMLAdapter struct {
	roomManager service.Manager
	Template    *template.Template
}

func NewGinHTMLAdapter(rm service.Manager) *ginHTMLAdapter {

	return &ginHTMLAdapter{rm, Html}
}

func (ga *ginHTMLAdapter) Stream(c *gin.Context) {
	roomid := c.Param("roomid")
	listener := ga.roomManager.OpenListener(roomid)
	defer ga.roomManager.CloseListener(roomid, listener)

	clientGone := c.Request.Context().Done()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(service.Message)
			if !ok {
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent("message", " "+serviceMsg.UserId+" → "+serviceMsg.Text)
			return true
		}
	})
}

func (ga *ginHTMLAdapter) GetRoom(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := fmt.Sprint(rand.Int31())
	c.HTML(http.StatusOK, "chat_room", gin.H{
		"roomid": roomid,
		"userid": userid,
	})
}

func (ga *ginHTMLAdapter) PostRoom(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := c.PostForm("user")
	message := c.PostForm("message")
	ga.roomManager.Submit(userid, roomid, message)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
	})
}

func (ga *ginHTMLAdapter) DeleteRoom(c *gin.Context) {
	roomid := c.Param("roomid")
	ga.roomManager.DeleteBroadcast(roomid)
}
