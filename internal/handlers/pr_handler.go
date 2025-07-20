package handlers

import (
	"strconv"

	"github.com/digaso/scalabit/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v57/github"
	"net/http"
)

func ListPRs(c *gin.Context) {
	ctx := c.Request.Context()

	owner := c.Param("owner")
	repo := c.Param("repo")

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	opts := &github.PullRequestListOptions{
		State: "open",
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	prs, _, err := client.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.CleanPRs(prs))
}
