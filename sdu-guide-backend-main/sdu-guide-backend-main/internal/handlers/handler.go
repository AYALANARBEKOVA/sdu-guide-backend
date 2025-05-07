package handlers

import (
	"sdu-guide/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
)

type Handler struct {
	Gin     *gin.Engine
	Service *services.Service
	Cache   *cache.Cache
}

func NewHandler(service *services.Service, cache *cache.Cache) *Handler {
	return &Handler{
		Gin:     gin.Default(),
		Service: service,
		Cache:   cache,
	}
}

func (h *Handler) Router() {
	h.Gin.POST("/sign-in", h.signIn())
	h.Gin.POST("/sign-up", h.signUp())
	h.Gin.GET("/xlsx/:hash", h.getXLSX)
	h.Gin.GET("/room/:id", h.getRoom())
	h.Gin.GET("/image/:hash", h.getImage)
	h.Gin.GET("/getAll-events", h.getAllEvents())
	h.Gin.GET("/event/:id", h.getEvent())
	h.Gin.GET("/getAll-events-currentMonth", h.getAllEventsForCallendar())
	h.Gin.GET("/schedule/:sef", h.getSchedule())
	h.Gin.GET("/translations", h.getLanguage())
	protectedGroup := h.Gin.Group("/")
	protectedGroup.Use(h.AuthRequired)
	{
		protectedGroup.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"Status": "Check"})
		})

		protectedGroup.GET("/logout", h.logout())
		protectedGroup.POST("/upload-XLSX", h.uploadXLSX())
		protectedGroup.POST("/upload-image", h.uploadImage())
		protectedGroup.POST("/create-room", h.createRoom())
		protectedGroup.PUT("/update-room", h.updateRoom())
		protectedGroup.GET("/profile", h.getProfile())
		protectedGroup.GET("/getAll-rooms", h.getAllRooms())
		protectedGroup.PUT("/update-user", h.updateUser())
		protectedGroup.DELETE("/delete-room/:id", h.deleteRoom())
		protectedGroup.POST("/create-event", h.createEvent())
		protectedGroup.PUT("/update-event", h.updateEvent())
		protectedGroup.DELETE("/delete-event/:id", h.deleteEvent())

	}
}
