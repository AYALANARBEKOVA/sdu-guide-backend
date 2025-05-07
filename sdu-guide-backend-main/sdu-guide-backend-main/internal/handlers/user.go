package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sdu-guide/internal/logger"
	"sdu-guide/internal/structures"
	"sdu-guide/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
)

func (h *Handler) signUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user *structures.UserRegister

		err := json.NewDecoder(c.Request.Body).Decode(&user)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}

		err = h.Service.UserService.Create(*user)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "User successfuly create"})
	}
}

func (h *Handler) signIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user structures.Login

		err := json.NewDecoder(c.Request.Body).Decode(&user)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		cookie, session, err := h.Service.UserService.LoginToSystem(user)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		err = h.Cache.Set(&cache.Item{Ctx: c.Request.Context(), Key: "session", Value: session, TTL: time.Minute * 60})
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		http.SetCookie(c.Writer, cookie)
		c.JSON(200, gin.H{"Status": "User successfuly signed"})
	}
}

func (h *Handler) logout() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Request.Cookie("Token")
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"Error": err.Error()})
			return
		}
		err = h.Service.Logout(cookie)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"Error": err.Error()})
			return
		}

		c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

		c.JSON(200, gin.H{"Status": "User successfuly logout from the system"})

	}
}

func (h *Handler) getProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session structures.Session

		if err := h.Cache.Get(context.Background(), "session", &session); err != nil {
			logger.Error.Println(err)
			c.JSON(401, gin.H{"error": "Not authorized"})
			return
		}

		user, err := h.Service.UserService.GetUserFromSession(session)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(401, gin.H{"error": "Not authorized"})
			return
		}
		c.JSON(200, gin.H{"data": user})
	}
}

func (h *Handler) updateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user *structures.User

		err := json.NewDecoder(c.Request.Body).Decode(&user)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": "Bad request"})
			return
		}
		fmt.Println("user:", user)
		err = h.Service.UserService.Upadte(*user)
		if err != nil {
			logger.Error.Println(err)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Bad request (%s)", err.Error())})
			return
		}

		c.JSON(200, gin.H{"Status": "User successfuly updated"})
	}
}

func (h *Handler) getLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {

		lang := c.Query("lang") // Получаем язык из параметра ?lang=
		if lang == "" {
			lang = "en" // Значение по умолчанию
		}

		translations, err := utils.LoadTranslations(lang)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, translations)
	}
}
