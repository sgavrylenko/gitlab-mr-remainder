package main

import (
	"fmt"
	"github.com/sgavrylenko/gitlab-mr-remander/internal/config"
	"os"
	"time"

	"github.com/xanzy/go-gitlab"
)

const mrPerPage = 100

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	git, err := gitlab.NewClient(appConfig.GitlabAuthToken, gitlab.WithBaseURL(appConfig.GitLabBaseURL))
	if err != nil {
		fmt.Println(err)
	}

	mrOpts := &gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: mrPerPage,
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
			if mergeRequest.ProjectID == appConfig.GitLabProjectID {
				fmt.Println(
					mergeRequest.Title, "|",
					mergeRequest.Description, "|",
					mergeRequest.DetailedMergeStatus, "|",
					fmt.Sprintf("is %.0f days old", time.Now().Sub(mergeRequest.UpdatedAt.UTC()).Hours()/24), "|",
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
