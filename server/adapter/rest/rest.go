package adapter

import (
	"io"
	"net/http"

	"github.com/MohammadBnei/realtime-chat/server/service"

	"github.com/gin-gonic/gin"
)

type GinAdapter interface {
	Stream(c *gin.Context)
	PostRoom(c *gin.Context)
	DeleteRoom(c *gin.Context)
}

type messageInput struct {
	UserId string `json:"userId"`
	Data   string `json:"data"`
}

type response struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

type ginAdapter struct {
	roomManager service.Manager
}

func NewGinAdapter(rm service.Manager) *ginAdapter {

	return &ginAdapter{rm}
}

// Stream godoc
// @Summary      Stream messages
// @Description  Stream messages from a room
// @Tags         chat
// @Produce      text/event-stream
// @Param        id   path      string  true  "Room ID"
// @Failure      400  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /stream/{id} [get]
func (ga *ginAdapter) Stream(c *gin.Context) {
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

// PostRoom godoc
// @Summary      Post to room
// @Description  Post a message to a room
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Room ID"
// @Param        messageInput   body      messageInput  true  "Message body"
// @Success      200  {object}  response
// @Failure      400  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /room/{id} [post]
func (ga *ginAdapter) PostRoom(c *gin.Context) {
	roomid := c.Param("roomid")
	var input messageInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	ga.roomManager.Submit(input.UserId, roomid, input.Data)

	c.JSON(http.StatusOK, &response{
		Success: true,
	})
}

// DeleteRoom godoc
// @Summary      Delete a room
// @Description  Delete the room
// @Tags         chat
// @Produce      json
// @Param        id   path      string  true  "Room ID"
// @Success      200  {object}  response
// @Failure      500  {object}  HTTPError
// @Router       /room/{id} [delete]
func (ga *ginAdapter) DeleteRoom(c *gin.Context) {
	roomid := c.Param("roomid")
	ga.roomManager.DeleteBroadcast(roomid)
	c.JSON(http.StatusOK, &response{
		Success: true,
	})
}
