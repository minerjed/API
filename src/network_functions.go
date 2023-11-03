package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func send_http_data(url string, data string) (string, error) {
	// create the http request
	fmt.Printf("2Entering send_http_data\n")
	fmt.Printf("url: %v\n", url)
	fmt.Printf("data: %v\n", data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", errors.New("failed to create HTTP request")
	}

	// set the request headers
	req.Header.Set("Content-Type", "application/json")

	// set the http client settings
	client := &http.Client{
		Timeout: time.Second * 2,
	}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("failed to send HTTP request")
	}
	defer resp.Body.Close()

	// get the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to read response body")
	}

	return string(body), nil
}
