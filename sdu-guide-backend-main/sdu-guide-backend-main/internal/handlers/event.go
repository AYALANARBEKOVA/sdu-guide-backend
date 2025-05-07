package handlers

import (
	"encoding/json"
	"fmt"
	"sdu-guide/internal/conv"
	"sdu-guide/internal/logger"
	"sdu-guide/internal/structures"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) createEvent() gin.HandlerFunc {
	return func(c *gin.Context) {

		var event *structures.Event

		err := json.NewDecoder(c.Request.Body).Decode(&event)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}
		if err := h.Service.EventService.CreateEvent(*event); err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "Room successfuly create"})
	}
}

func (h *Handler) getEvent() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		event, err := h.Service.EventService.GetEvent(conv.Int64(id))
		if err != nil {
			logger.Error.Println(err)
			c.JSON(500, gin.H{"error": "Can't get event"})
			return
		}

		c.JSON(200, gin.H{"data": event})

	}
}

func (h *Handler) updateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {

		var event structures.Event

		err := json.NewDecoder(c.Request.Body).Decode(&event)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}

		if err := h.Service.EventService.UpdateEvent(event); err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "Event successfuly updated"})
	}
}

func (h *Handler) getAllEvents() gin.HandlerFunc {
	return func(c *gin.Context) {

		limit := c.Query("limit")
		ended := c.Query("withEnded")
		today := c.Query("today")

		filter := structures.Filter{
			Request: bson.M{},
		}

		if len(ended) > 0 {
			if ended == "false" {
				filter.Request["ended"] = false
			}
		} else {
			filter.Request["ended"] = false
		}
		if limit != "" {
			if parsedLimit, err := strconv.ParseInt(limit, 10, 64); err == nil {
				filter.Limit = parsedLimit
			} else {
				c.JSON(400, gin.H{"error": "Invalid limit value"})
				return
			}
		}

		if today != "" {
			if today == "true" {
				now := time.Now()
				startOfMonth := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

				filter.Request["date"] = bson.M{
					"$gte": startOfMonth,
					"$lt":  now.Add(24 * time.Hour),
				}
			}
		}

		result, err := h.Service.EventService.GetAll(filter)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, bson.M{"data": result})
	}
}

func (h *Handler) deleteEvent() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		if conv.Uint64(id) == 0 {
			c.JSON(400, gin.H{"error": "Wrong id format"})
			return
		}
		if err := h.Service.EventService.Delete(conv.Uint64(id)); err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "Event successfuly deleted"})
	}
}

func (h *Handler) getAllEventsForCallendar() gin.HandlerFunc {
	return func(c *gin.Context) {

		now := time.Now()

		// Определяем начало и конец текущего месяца
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0) // Первый день следующего месяца

		filter := structures.Filter{
			Request: bson.M{
				"date": bson.M{
					"$gte": startOfMonth, // Даты с начала текущего месяца
					"$lt":  endOfMonth,   // Даты до следующего месяца
				},
			},
		}

		result, err := h.Service.EventService.GetAll(filter)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, bson.M{"data": result})
	}
}
