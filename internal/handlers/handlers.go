package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/cloud-server/grpc"
	"github.com/zarinit-routers/connector-rpc/gen/connector"
)

func GetClients() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		response, err := grpc.NodesService.NodesByGroup(
			ctx,
			&connector.NodesByGroupRequest{
				GroupId: dummyId(),
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func dummyId() string {
	return "00000000-0000-0000-0000-00000000000"
}
