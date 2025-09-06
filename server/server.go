package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zarinit-routers/cloud-server/internal/handlers"
)

func init() {
	// gin.SetMode(gin.ReleaseMode)
}

type Server struct {
	engine *gin.Engine
}

func (s *Server) Start() error {
	addr := getAddr()
	return s.engine.Run(addr)
}

func New() *Server {
	engine := gin.Default()

	handlers.SetupNodeCommands(engine.Group("/api/cmd"))

	return &Server{
		engine: engine,
	}
}

func getAddr() string {
	return ":9090"
}

func setupRoutes(r *gin.RouterGroup) {
}
