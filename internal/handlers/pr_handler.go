package handlers

import (
	"context"

	"github.com/digaso/scalabit/internal/utils"

	"github.com/gin-gonic/gin"

	"github.com/google/go-github/v57/github"

	"net/http"

	"strconv"
)

func ListPRs(c *gin.Context) {

	ctx := context.Background()

	repo := c.Param("repo")

	owner := c.Param("owner")

	limitStr := c.DefaultQuery("limit", "10") // Default limit to 10 if not provided

	limit, err := strconv.Atoi(limitStr)

	if err != nil || limit <= 0 {

		limit = 10

	}

	opts := &github.PullRequestListOptions{

		ListOptions: github.ListOptions{PerPage: limit},

		State: "open",
	}

	prs, _, err := client.PullRequests.List(ctx, owner, repo, opts)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return

	}

	c.JSON(http.StatusOK, utils.CleanPRs(prs))

}
