package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func buildRS(filename string, replacement string) *repoStruct {
	testFile, err := os.Open("./test-data/" + filename)
	check(err)
	defer testFile.Close()
	b, err := ioutil.ReadAll(testFile)
	check(err)
	results := fmt.Sprintf(string(b), replacement)
	bytes := []byte(results)
	rs := repoStruct{}
	parseProtectionDetails(&rs, bytes)
	return &rs
}

func TestParseProtectionDetails(t *testing.T) {
	// protected is true, no context, include admins
	prot :=
		`"protection": {
			"enabled": true,
			"required_status_checks": {
			"enforcement_level": "everyone",
			"contexts": []
		}
    }`

	rs := buildRS("protected_response.json", prot)
	assert.Equal(t, true, rs.Protected)
	assert.Equal(t, false, rs.ProtectedWithStatusCheck)
}

func TestParsePrtectionDetailsProtectedWithStatusCheck(t *testing.T) {
	// protected is true, context is build, include admins
	prot :=
		`"protection": {
			"enabled": true,
			"required_status_checks": {
			"enforcement_level": "everyone",
			"contexts": ["build"]
		}
    }`

	rs := buildRS("protected_response.json", prot)

	assert.Equal(t, false, rs.Protected)
	assert.Equal(t, true, rs.ProtectedWithStatusCheck)
}

func TestParsePrtectionDetailsNotProtected(t *testing.T) {
	// protected is false
	prot :=
		`"protection": {
			"enabled": false,
			"required_status_checks": {
			"enforcement_level": "everyone",
			"contexts": []
		}
    }`

	rs := buildRS("protected_response.json", prot)

	assert.Equal(t, false, rs.Protected)
	assert.Equal(t, false, rs.ProtectedWithStatusCheck)
}

func TestParsePrtectionDetailsNoAdmins(t *testing.T) {
	// not include admins
	prot :=
		`"protection": {
			"enabled": true,
			"required_status_checks": {
			"enforcement_level": "non_admins",
			"contexts": []
		}
    }`

	rs := buildRS("protected_response.json", prot)

	assert.Equal(t, false, rs.Protected)
	assert.Equal(t, false, rs.ProtectedWithStatusCheck)
}

func TestParseHooksEmpty(t *testing.T) {
	rs := repoStruct{}
	var hooks []github.Hook
	parseHookData(&rs, hooks, "")
	assert.Equal(t, false, rs.CIHooks)
	assert.Equal(t, false, rs.CITestHooks)
}

func TestParseHooksCI(t *testing.T) {
	url := "https://%s.test.com/"
	rs := repoStruct{}
	hooks := []github.Hook{
		*getProdPushWH(),
		*getProdPRWH(),
	}

	parseHookData(&rs, hooks, url)

	assert.Equal(t, true, rs.CIHooks)
	assert.Equal(t, false, rs.CITestHooks)
}

func getProdPushWH() *github.Hook {
	return &github.Hook{
		Events: []string{"push"},
		Config: map[string]interface{}{
			"url":          "https://ci-proxy.test.com/github-webhook/",
			"content_type": "form",
			"insecure_ssl": "0",
		},
	}
}

func getProdPRWH() *github.Hook {
	return &github.Hook{
		Events: []string{"issue_comment", "pull_request"},
		Config: map[string]interface{}{
			"url":          "https://ci-proxy.test.com/ghprbhook/",
			"content_type": "form",
			"insecure_ssl": "0",
		},
	}
}

func getTestPushWH() *github.Hook {
	return &github.Hook{
		Events: []string{"push"},
		Config: map[string]interface{}{
			"url":          "https://ci-test-proxy.test.com/github-webhook/",
			"content_type": "form",
			"insecure_ssl": "0",
		},
	}
}

func getTestPRWH() *github.Hook {
	return &github.Hook{
		Events: []string{"issue_comment", "pull_request"},
		Config: map[string]interface{}{
			"url":          "https://ci-test-proxy.test.com/ghprbhook/",
			"content_type": "form",
			"insecure_ssl": "0",
		},
	}
}
