package infra

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/sirupsen/logrus"
)

// GinServer is the struct gathering all the server details
type GinServer struct {
	host   string
	port   int
	Router *gin.Engine
}

// NewServer creates the gin Server
func NewServer(host string, port int, mode, allowedOrigins string) GinServer {
	s := GinServer{}
	s.port = port
	s.host = host

	s.Router = gin.New()
	s.setMode(mode)
	s.Router.Use(gin.Recovery())
	s.setCORS(allowedOrigins)

	return s
}

func (s *GinServer) setMode(mode string) {
	switch mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		logrus.WithField("mode", mode).Warn("Unknown gin mode, fallback to release")
		gin.SetMode(gin.ReleaseMode)
	}
}

// setCORS is a helper to set current engine CORS
func (s *GinServer) setCORS(allowedOrigins string) {
	s.Router.Use(cors.Middleware(cors.Config{
		Origins:         allowedOrigins,
		Methods:         "GET, PUT, POST, DELETE, OPTION, PATCH",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
}

// Start tells the router to start listening
func (s GinServer) Start() {
	if err := s.Router.Run(fmt.Sprintf("%s:%d", s.host, s.port)); err != nil {
		logrus.WithError(err).Fatal("Couldn't start router")
	}
}
