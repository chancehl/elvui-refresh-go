package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v64/github"
)

// getLatestTaggedVersion fetches the latest tagged version of ElvUI from the GitHub repository.
func GetLatestTaggedVersion() (*github.RepositoryTag, error) {
	client := github.NewClient(nil)
	options := &github.ListOptions{}

	tags, _, err := client.Repositories.ListTags(context.Background(), TukUI, ElvUI, options)
	if err != nil {
		fmt.Printf("Could not fetch ElvUI repo: %v", err)
		return nil, err
	}

	return tags[0], nil
}
