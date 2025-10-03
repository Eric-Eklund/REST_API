package routes

import (
	"REST_API/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerEvent(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Event not found"})
		return
	}

	err = event.Register(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Event could not be registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event registered successfully"})
}

func unregisterEvent(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Event not found"})
		return
	}

	err = event.Unregister(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Event could not be unregistered"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event unregistered successfully"})
}
