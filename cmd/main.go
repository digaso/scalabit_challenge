package main

import (
	"github.com/digaso/scalabit/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN not set in environment variables")
	}

	handlers.InitGitHubClient(token)

	router := gin.Default()

	router.GET("/repos/:owner", handlers.ListRepos)
	router.POST("/repos", handlers.CreateRepo)
	router.DELETE("/repos", handlers.DeleteRepo)
	router.GET("/repos/prs/:owner/:repo", handlers.ListPRs)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
