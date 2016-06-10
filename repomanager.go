package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Will limit the count of results for testing purposes
var debug = true

type repoStruct struct {
	Name      *string
	Protected bool
	Contexts  []string
}

type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		Sha    string `json:"sha"`
		Commit struct {
			Author struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Committer struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"committer"`
			Message string `json:"message"`
			Tree    struct {
				Sha string `json:"sha"`
				URL string `json:"url"`
			} `json:"tree"`
			URL          string `json:"url"`
			CommentCount int    `json:"comment_count"`
		} `json:"commit"`
		URL         string `json:"url"`
		HTMLURL     string `json:"html_url"`
		CommentsURL string `json:"comments_url"`
		Author      struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"author"`
		Committer struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"committer"`
		Parents []struct {
			Sha     string `json:"sha"`
			URL     string `json:"url"`
			HTMLURL string `json:"html_url"`
		} `json:"parents"`
	} `json:"commit"`
	Links struct {
		Self string `json:"self"`
		HTML string `json:"html"`
	} `json:"_links"`
	Protection struct {
		Enabled              bool `json:"enabled"`
		RequiredStatusChecks struct {
			EnforcementLevel string   `json:"enforcement_level"`
			Contexts         []string `json:"contexts"`
		} `json:"required_status_checks"`
	} `json:"protection"`
}

// GetAllRepos will return a list of repos for a Github Organization
func getAllRepos(orgname string) []repoStruct {
	client := getClient()

	var count int
	if debug {
		count = 5
	} else {
		count = 100
	}
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: count},
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

		if debug {
			break
		}
	}
	log.Printf("Found %d repo(s) for the organization %s", len(allRepos), orgname)

	var repos []repoStruct
	repos = make([]repoStruct, len(allRepos))
	for i, repo := range allRepos {
		rs := repoStruct{}
		rs.Name = repo.Name
		addProtectedDetails(&rs, orgname, repo)
		repos[i] = rs
	}
	return repos
}

func addProtectedDetails(rs *repoStruct, orgname string, repo github.Repository) {
	url := "https://api.github.com/repos/%s/%s/branches/master"
	fullURL := fmt.Sprintf(url, orgname, *repo.Name)
	log.Printf("Connecting to URL: %s", fullURL)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("Authorization", "token "+*_config.githubkey)
	req.Header.Set("accept", "application/vnd.github.loki-preview+json")
	resp, err := client.Do(req)

	check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	var s = new(Branch)
	err = json.Unmarshal(body, &s)
	check(err)

	rs.Protected = s.Protection.Enabled
	rs.Contexts = s.Protection.RequiredStatusChecks.Contexts
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
