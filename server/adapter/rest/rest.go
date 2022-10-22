package adapter

import (
	"io"
	"net/http"
	"realtime-chat/service"

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
// @Produce      octet-stream
// @Param        id   path      string  "roomie"  "Room ID"
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /rooms/{id} [get]
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
			c.SSEvent("message", message)
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
// @Param        id   path      string  "roomie"  "Room ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /rooms/{id} [post]
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
// @Param        id   path      string  "roomie"  "Room ID"
// @Success      200  {object}  Response
// @Failure      500  {object}  httputil.HTTPError
// @Router       /rooms/{id} [delete]
func (ga *ginAdapter) DeleteRoom(c *gin.Context) {
	roomid := c.Param("roomid")
	ga.roomManager.DeleteBroadcast(roomid)
	c.JSON(http.StatusOK, &response{
		Success: true,
	})
}
