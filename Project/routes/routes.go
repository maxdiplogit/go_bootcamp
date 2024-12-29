package routes

import (
	"example.com/app/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	authenticatedRoutes := server.Group("/")

	authenticatedRoutes.Use(middlewares.Authenticate)
	authenticatedRoutes.POST("/events", createEvent)
	authenticatedRoutes.PUT("/events/:id", updateEvent)
	authenticatedRoutes.DELETE("/events/:id", deleteEvent)
	authenticatedRoutes.POST("/events/:id/register", registerForEvent)
	authenticatedRoutes.DELETE("/events/:id/register", cancelRegistration)

	server.POST("/register", signup)
	server.POST("/login", login)
}
