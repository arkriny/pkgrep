package voidlinux

import (
	"fmt"
	"net/http"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/void-linux/void-packages/contents/srcpkgs/%s", query)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK
	return ok, nil
}
