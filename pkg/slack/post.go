package slack

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Post(token, channel, message string) error {
	values := url.Values{}
	values.Set("token", token)
	values.Set("channel", channel)
	values.Set("text", message)

	request, err := http.NewRequest(
		"POST",
		"https://slack.com/api/chat.postMessage",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer response.Body.Close()

	return nil
}
