package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type issue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

func PostIssue(token, owner, repositoryName, title, body string, labels []string) error {
	json, err := json.Marshal(issue{
		Title:  title,
		Body:   body,
		Labels: labels,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repositoryName),
		bytes.NewBuffer(json),
	)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
