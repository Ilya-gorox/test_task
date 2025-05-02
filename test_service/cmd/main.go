package main

import (
	"log"
	"os"

	"test_service/internal/handler"
	"test_service/internal/repository"
	"test_service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	router := gin.Default()

	router.POST("/users", handler.CreateUser)
	router.GET("/user/:id", handler.GetUser)
	router.PATCH("/user/:id", handler.UpdateUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
