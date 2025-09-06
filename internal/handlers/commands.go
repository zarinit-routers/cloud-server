package handlers

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/cloud-server/queue"
	cmd "github.com/zarinit-routers/commands"
)

var commands = []string{
	cmd.CommandTimezoneGet.String(),
	cmd.CommandTimezoneSet.String(),
}

func SetupNodeCommands(r *gin.RouterGroup) {

	for _, cmd := range commands {
		r.POST("/"+cmd, func(c *gin.Context) {

			var req queue.Request

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, ResponseErr(err))
				return
			}

			req.Command = cmd

			response, err := queue.SendRequest(&req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ResponseErr(err))
				return
			}

			c.JSON(http.StatusOK, response)
		})
	}
}

func ResponseErr(err error) queue.Response {
	log.Error("Error sending request", "error", err)
	return queue.Response{
		RequestError: fmt.Sprintf("pre-command error: bad request %s", err),
	}

}
