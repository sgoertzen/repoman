package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

type RepoStruct struct {
	Name *string
	Protected bool
}


// GetAllRepos will return a list of repos for a Github Organization
func getAllRepos(orgname string) []RepoStruct {
	client := getClient()

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(orgname, opt)
		if err != nil {
			return nil
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	log.Printf("Found %d repo(s) for the organization %s", len(allRepos), orgname)

	var repos []RepoStruct
	repos = make([]RepoStruct, len(allRepos))
	for i, repo := range allRepos {
		rs := RepoStruct{}
		rs.Name = repo.Name

		//TODO: Make call to Github api with beta header and get the value of protected
		rs.Protected = getProtectedStatus(repo)

		repos[i] = rs
	} 
	return repos
}

func getProtectedStatus(github.Repository) bool {
	return false
}

func getClient() *github.Client {
	var tc *http.Client
	envToken := *_config.githubkey
	if len(envToken) > 0 {
		token := oauth2.Token{AccessToken: envToken}
		ts := oauth2.StaticTokenSource(&token)
		tc = oauth2.NewClient(oauth2.NoContext, ts)
	}
	client := github.NewClient(tc)
	return client
}