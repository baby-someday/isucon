package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type issue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type PostIssueResponse struct {
	URL string `json:"url"`
}

func (p *PostIssueResponse) GetID() (int64, error) {
	componensts := strings.Split(p.URL, "/")
	idString := componensts[len(componensts)-1]
	return strconv.ParseInt(idString, 10, 64)
}

func PostIssue(token, owner, repositoryName, title, body string, labels []string) (PostIssueResponse, error) {
	requestBody, err := json.Marshal(issue{
		Title:  title,
		Body:   body,
		Labels: labels,
	})
	if err != nil {
		return PostIssueResponse{}, err
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repositoryName),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return PostIssueResponse{}, err
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return PostIssueResponse{}, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return PostIssueResponse{}, err
	}

	postIssueResponse := PostIssueResponse{}
	err = json.Unmarshal(responseBytes, &postIssueResponse)
	if err != nil {
		return PostIssueResponse{}, err
	}

	return postIssueResponse, nil
}
