package main

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"
	"strconv"
	"time"
)

func main() {
	GitlabAuthToken, ok := os.LookupEnv("MR_REMAINDER_TOKEN")
	if !ok || GitlabAuthToken == "" {
		fmt.Println("MR_REMAINDER_TOKEN environment variable not set")
		os.Exit(1)
	}
	GitLabBaseURL, ok := os.LookupEnv("CI_SERVER_URL")
	if !ok || GitLabBaseURL == "" {
		fmt.Println("CI_SERVER_URL environment variable not set")
		os.Exit(1)
	}

	gitLabProjectID, ok := os.LookupEnv("CI_PROJECT_ID")
	if !ok || gitLabProjectID == "" {
		fmt.Println("CI_PROJECT_ID environment variable not set")
	}
	GitLabProjectID, err := strconv.Atoi(gitLabProjectID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	git, err := gitlab.NewClient(GitlabAuthToken, gitlab.WithBaseURL(GitLabBaseURL))
	if err != nil {
		fmt.Println(err)
	}

	mrOpts := &gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Scope: gitlab.Ptr("all"),
		State: gitlab.Ptr("opened"),
		WIP:   gitlab.Ptr("no"),
	}

	for {
		mrs, resp, err := git.MergeRequests.ListMergeRequests(mrOpts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, mergeRequest := range mrs {
			if mergeRequest.ProjectID == GitLabProjectID {
				fmt.Println(
					mergeRequest.Title, "|",
					mergeRequest.Description, "|",
					mergeRequest.DetailedMergeStatus, "|",
					fmt.Sprintf("%.0f days old", time.Now().Sub(mergeRequest.UpdatedAt.UTC()).Hours()/24), "|",
					mergeRequest.WebURL,
				)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		mrOpts.Page = resp.NextPage
	}
}
