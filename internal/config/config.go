package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	GitlabAuthToken string
	GitLabBaseURL   string
	GitLabProjectID int
}

func NewAppConfig() (*AppConfig, error) {
	GitlabAuthToken, ok := os.LookupEnv("MR_REMAINDER_TOKEN")
	if !ok || GitlabAuthToken == "" {
		return nil, fmt.Errorf("MR_REMAINDER_TOKEN environment variable not set")
	}

	GitLabBaseURL, ok := os.LookupEnv("CI_SERVER_URL")
	if !ok || GitLabBaseURL == "" {
		return nil, fmt.Errorf("CI_SERVER_URL environment variable not set")
	}

	gitLabProjectID, ok := os.LookupEnv("CI_PROJECT_ID")
	if !ok || gitLabProjectID == "" {
		return nil, fmt.Errorf("CI_PROJECT_ID environment variable not set")
	}
	GitLabProjectID, err := strconv.Atoi(gitLabProjectID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &AppConfig{
		GitlabAuthToken: GitlabAuthToken,
		GitLabBaseURL:   GitLabBaseURL,
		GitLabProjectID: GitLabProjectID,
	}, nil
}
