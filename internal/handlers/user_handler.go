package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/digaso/scalabit/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v57/github"
)

var client *github.Client

type bearerAuthTransport struct {
	token     string
	transport http.RoundTripper
}

// Problems with authentication made me add "Bearer" prefix

func (bat *bearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+bat.token)
	return bat.transport.RoundTrip(req)
}

func InitGitHubClient(token string) {
	if token == "" {
		log.Fatal("GitHub token cannot be empty")
	}

	transport := &bearerAuthTransport{
		token:     token,
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

func DeleteRepo(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	repo := req.Name

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching authenticated user: " + err.Error()})
		return
	}
	// Ensure the authenticated user is the owner of the repository
	owner := *user.Login

	_, err = client.Repositories.Delete(ctx, owner, repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
