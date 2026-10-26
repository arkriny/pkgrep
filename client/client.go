package client

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
)

// getCheckOK is a helper function that sends request to URL and checks if
// response status code is 200 (OK).
func getCheckOK(client *http.Client, url string) (bool, error) {
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK
	return ok, nil
}

// getBodyContains is a helper function that sends request to URL and checks if
// response status code is 200 (OK) and response body contains specified substr.
func getBodyContains(client *http.Client, url, substr string) (bool, error) {
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	substrBytes := []byte(substr)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.Contains(line, substrBytes) {
			return true, nil
		}
	}
	return false, scanner.Err()
}
