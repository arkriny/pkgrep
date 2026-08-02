package httputil

import (
	"bufio"
	"errors"
	"net/http"
	"strings"
)

// GetCheckOK is a helper function that sends request to URL and checks if
// response status code is 200 (OK).
func GetCheckOK(client *http.Client, url string) (bool, error) {
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK
	return ok, nil
}

// GetBodyContains is a helper function that sends request to URL and checks if
// response status code is 200 (OK) and response body contains specified substr.
func GetBodyContains(client *http.Client, url, substr string) (bool, error) {
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, substr) {
			return true, nil
		}
	}
	return false, scanner.Err()
}
