package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Commit struct {
	Sha1 string `json:"sha"`
}

func (c *Commit) GetShortSha1() string {
	return c.Sha1[:8]
}

func GetCommit(token, owner, repository, branch string) (Commit, error) {
	request, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repository, branch),
		strings.NewReader(url.Values{}.Encode()),
	)
	if err != nil {
		return Commit{}, err
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return Commit{}, err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Commit{}, err
	}

	commit := Commit{}
	err = json.Unmarshal(bytes, &commit)
	if err != nil {
		return Commit{}, err
	}

	return commit, nil
}
