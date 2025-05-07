package main

import (
	"sdu-guide/internal/ai"
	"sdu-guide/internal/handlers"
	"sdu-guide/internal/repositories"
	"sdu-guide/internal/services"
	"time"

	"github.com/gin-contrib/cors"
)

func main() {
	db, err := repositories.NewDB()
	if err != nil {
		panic(err)
	}
	ai := ai.NewAI(db)
	repo := repositories.NewRepository(db, &ai)
	service := services.NewService(repo)
	cache, err := handlers.NewCache()
	if err != nil {
		panic(err)
	}
	handler := handlers.NewHandler(service, cache)

	config := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Access-Control-Allow-Origin", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	handler.Gin.Use(config)
	handler.Router()
	err = handler.Gin.Run("localhost:8000")
	if err != nil {
		panic(err)
	}
}
