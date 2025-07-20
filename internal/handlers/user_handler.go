package handlers

import (
	"context"

	"github.com/digaso/scalabit/internal/utils"

	"github.com/gin-gonic/gin"

	"github.com/google/go-github/v57/github"

	"github.com/joho/godotenv"

	"log"

	"net/http"

	"os"
)

var client *github.Client

type bearerAuthTransport struct {
	token string

	transport http.RoundTripper
}

func (bat *bearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	req.Header.Set("Authorization", "Bearer "+bat.token)

	return bat.transport.RoundTrip(req)

}

func init() {

	err := godotenv.Load()

	if err != nil {

		log.Print("No .env file found or error loading .env")

	}

	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {

		log.Fatal("GITHUB_TOKEN not available in environment variables")

	}

	transport := &bearerAuthTransport{

		token: token,

		transport: http.DefaultTransport,
	}

	client = github.NewClient(&http.Client{Transport: transport})

	if client == nil {

		log.Fatal("Failed to create GitHub client")

	}

	log.Print("GitHub client created successfully")

}

func ListRepos(c *gin.Context) {

	ctx := context.Background()

	owner := c.Param("owner")

	repos, _, err := client.Repositories.ListByUser(ctx, owner, nil)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return

	}

	c.JSON(http.StatusOK, utils.CleanRepos(repos))

}

func CreateRepo(c *gin.Context) {

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})

		return

	}

	ctx := context.Background()

	repo := &github.Repository{

		Name: github.String(req.Name),
	}

	createdRepo, _, err := client.Repositories.Create(ctx, "", repo)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return

	}

	c.JSON(http.StatusCreated, createdRepo)

}

// DeleteRepo deletes a repository by name

func DeleteRepo(c *gin.Context) {

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})

		return

	}

	//I do this so we don't have to pass the owner in the request body

	//This assumes the authenticated user is the owner of the repo because we can only delete repos we own

	ctx := context.Background()

	repo := string(req.Name)

	user, _, err := client.Users.Get(context.Background(), "")

	if err != nil {

		log.Fatal("Error fetching authenticated user:", err)

	}

	owner := *user.Login

	_, err = client.Repositories.Delete(ctx, owner, repo)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return

	}

	c.JSON(http.StatusNoContent, nil)

}
