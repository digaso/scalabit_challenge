package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v57/github"
)

// CleanRepos formats the repository data for the response

func CleanRepos(repos []*github.Repository) []gin.H {
	var cleanedRepos = []gin.H{}
	for _, repo := range repos {
		cleanedRepos = append(cleanedRepos, gin.H{
			"name":        repo.GetName(),
			"description": repo.GetDescription(),
			"owner":       repo.Owner.GetLogin(),
			"url":         repo.GetHTMLURL(),
		})
	}
	return cleanedRepos
}

// CleanPRs formats the pull request data for the response

func CleanPRs(prs []*github.PullRequest) []gin.H {
	var cleanedPRs = []gin.H{}
	for _, pr := range prs {
		cleanedPRs = append(cleanedPRs, gin.H{
			"title":  pr.Title,
			"user":   pr.User.GetName(),
			"number": pr.GetNumber(),
		})
	}
	return cleanedPRs
}
