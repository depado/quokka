package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Interface is the interface that has to be implemented in order to handle
// every functionality of the router
type Interface interface {
	Health(c *gin.Context)
}

// New returns a new router
func New() Interface {
	return Impl{}
}

// Impl is the implementation
type Impl struct{}

// Health is the implementation of health check
func (r Impl) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": true})
}

//SetRoutes update router with routes
func SetRoutes(router *gin.Engine, r Interface) {
	router.GET("/health", r.Health)
}
