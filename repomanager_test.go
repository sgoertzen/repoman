package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestFileContents(filename string) []byte {
	testFile, err := os.Open("./test-data/" + filename)
	check(err)
	defer testFile.Close()
	b, err := ioutil.ReadAll(testFile)
	check(err)
	return b
}

func TestParsePrtectionDetails(t *testing.T) {
	resp := getTestFileContents("protected_response.json")
	rs := repoStruct{}
	parseProtectionDetails(&rs, resp)
	assert.Equal(t, true, rs.Protected)
}
