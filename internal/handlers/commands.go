package handlers

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/cloud-server/queue"
)

func SetupNodeCommands(r *gin.RouterGroup) {

	r.POST("/", func(c *gin.Context) {
		var req queue.Request

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Failed bind request body", "error", err)
			c.JSON(http.StatusBadRequest, ResponseErr(err))
			return
		}

		if req.Command == "" {
			log.Error("No command specified")
			c.JSON(http.StatusBadRequest, ResponseErr(fmt.Errorf("No command specified")))
			return
		}
		if req.NodeId == "" {
			log.Error("No node ID specified")
			c.JSON(http.StatusBadRequest, ResponseErr(fmt.Errorf("No node ID specified")))
			return
		}

		response, err := queue.SendRequest(&req)
		if err != nil {
			log.Error("Failed send/get request", "error", err)
			c.JSON(http.StatusInternalServerError, ResponseErr(err))
			return
		}

		c.JSON(http.StatusOK, response)
	})
}

func ResponseErr(err error) queue.Response {
	log.Error("Error sending request", "error", err)
	return queue.Response{
		RequestError: fmt.Sprintf("pre-command error: bad request %s", err),
	}

}
