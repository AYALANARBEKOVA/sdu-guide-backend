package handlers

import (
	"sdu-guide/internal/logger"
	"sdu-guide/internal/structures"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthRequired(c *gin.Context) {

	cookie, err := c.Request.Cookie("Token")
	if err != nil {
		logger.Error.Println(err)
		c.JSON(401, gin.H{"error": "Unauthorized", "status": 0})
		c.Abort()
		return
	}
	var session structures.Session
	err = h.Cache.Get(c.Request.Context(), "session", &session)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(401, gin.H{"error": "Unauthorized", "status": 0})
		c.Abort()
		return
	}

	if cookie.Value != session.Token || session.ExpiredDate.Before(time.Now()) {
		logger.Error.Println(err)
		h.Cache.Delete(c.Request.Context(), "session")
		c.JSON(401, gin.H{"error": "Unauthorized", "status": 0})
		c.Abort()
		return
	}

	c.Next()
}
