package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/digaso/scalabit/internal/handlers"
	"github.com/gin-gonic/gin"
	//"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	// Load .env once (locally development)
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Println("No .env file found or error loading .env")
	// }

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN not set in environment variables")
	}

	handlers.InitGitHubClient(token)

	// Setup router
	router = gin.Default()
	router.POST("/repos", handlers.CreateRepo)
	router.DELETE("/repos", handlers.DeleteRepo)

	code := m.Run()
	os.Exit(code)
}

func TestCreateRepo(t *testing.T) {
	repoName := "scalabit-test-repo"

	body := map[string]string{
		"name": repoName,
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/repos", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var createdRepo map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &createdRepo)
	assert.NoError(t, err)

	assert.Equal(t, repoName, createdRepo["name"], "repo name should match")
}

func TestDeleteRepo(t *testing.T) {
	repoName := "scalabit-test-repo"

	body := map[string]string{
		"name": repoName,
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("DELETE", "/repos", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}
