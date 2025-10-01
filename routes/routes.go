package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	// Events
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEventByID)

	server.POST("/events", createEvent)

	server.PUT("/events/:id", updateEvents)

	server.DELETE("/events/:id", deleteEvent)

	// Users
	server.POST("/signup", signup)
	server.POST("/login", login)
}
