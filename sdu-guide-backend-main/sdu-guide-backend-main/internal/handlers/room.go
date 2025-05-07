package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sdu-guide/internal/conv"
	"sdu-guide/internal/logger"
	"sdu-guide/internal/structures"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) uploadXLSX() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("xlsx")
		if err != nil {
			logger.Error.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file"})
			return
		}

		hash, err := h.Service.FileServices.StoreFile(file, header)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"hash": hash})
	}
}

func (h *Handler) uploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			logger.Error.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file"})
			return
		}

		hash, err := h.Service.FileServices.StoreFile(file, header)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get the file"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"hash": hash})
	}
}

func (h *Handler) getXLSX(c *gin.Context) {

	hash := c.Param("hash")
	xlsx, err := h.Service.FileServices.GetFilebyHash(hash)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get xlsx: " + err.Error()})
		return
	}

	data, err := os.ReadFile(xlsx.Path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Error.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			logger.Error.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file: " + err.Error()})
		}
		return
	}

	// Устанавливаем корректный Content-Type
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+xlsx.Name) // Чтобы при скачивании был нормальный filename
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

}

func (h *Handler) getImage(c *gin.Context) {

	hash := c.Param("hash")
	image, err := h.Service.FileServices.GetFilebyHash(hash)
	if err != nil {
		logger.Error.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get xlsx: " + err.Error()})
		return
	}

	data, err := os.ReadFile(image.Path)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file: " + err.Error()})
		}
		return
	}

	contentType := http.DetectContentType(data)
	c.Data(http.StatusOK, contentType, data)

}

func (h *Handler) createRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		var room *structures.Room

		err := json.NewDecoder(c.Request.Body).Decode(&room)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}
		if err := h.Service.RoomService.CreateRoom(*room); err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "Room successfuly create"})
	}
}

func (h *Handler) getRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		room, err := h.Service.RoomService.GetRoom(conv.Int64(id))
		if err != nil {
			logger.Error.Println(err)
			c.JSON(500, gin.H{"error": "Can't get room"})
			return
		}

		c.JSON(200, gin.H{"data": room})

	}
}

func (h *Handler) updateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		var room structures.Room

		err := json.NewDecoder(c.Request.Body).Decode(&room)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}
		logger.Info.Println("ROOM:", room)
		if err := h.Service.RoomService.UpdateRoom(room); err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "Room successfuly updated"})
	}
}

func (h *Handler) getAllRooms() gin.HandlerFunc {
	return func(c *gin.Context) {

		block := c.Query("block")
		number := conv.Int64(c.Query("number"))

		filter := bson.M{}
		if len(block) != 0 {
			filter["block"] = block
		}
		if number > 0 {
			filter["number"] = number
		}

		result, err := h.Service.RoomService.GetAll(filter)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, bson.M{"data": result})
	}
}

func (h *Handler) deleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		if conv.Uint64(id) == 0 {
			c.JSON(400, gin.H{"error": "Wrong id format"})
			return
		}
		if err := h.Service.RoomService.Delete(conv.Uint64(id)); err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "Room successfuly deleted"})
	}
}

func (h *Handler) getSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {

		sef := c.Param("sef")
		if sef == "" {
			c.JSON(500, gin.H{"error": "SEF can't be empty"})
			return
		}
		xlsx, err := h.Service.RoomService.GetScheduleBySEF(sef)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(500, gin.H{"error": "Can't get event"})
			return
		}

		data, err := os.ReadFile(xlsx.Path)
		if err != nil {
			if os.IsNotExist(err) {
				logger.Error.Println(err)
				c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			} else {
				logger.Error.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file: " + err.Error()})
			}
			return
		}

		// Устанавливаем корректный Content-Type
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename="+xlsx.Name) // Чтобы при скачивании был нормальный filename
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	}
}
