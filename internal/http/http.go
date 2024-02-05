package http

import (
	"esensi-test/internal/app/auth"
	"esensi-test/internal/app/customer"
	"esensi-test/internal/app/item"
	"esensi-test/internal/app/user"
	"esensi-test/internal/factory"
	"esensi-test/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Here we define route function for user Handlers that accepts gin.Engine and factory parameters
func NewHttp(g *gin.Engine, f *factory.Factory) {
	// Here we use logger middleware before the actual API to catch any api call from clients
	g.Use(gin.Logger())
	// Here we use the recovery middleware to catch a panic, if panic occurs recover the application witohut shutting it off
	g.Use(gin.Recovery())

	g.Use(middleware.CORSMiddleware())

	// Here we define a router group
	v1 := g.Group("/api/v1")
	// Here we register the route from user handler
	auth.NewHandler(f).Router(v1.Group("/auth"))
	user.NewHandler(f).Router(v1.Group("/user"))
	item.NewHandler(f).Router(v1.Group("/item"))
	customer.NewHandler(f).Router(v1.Group("/customer"))

}
