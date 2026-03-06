package homebrew

import (
	"fmt"
	"net/http"
)

func Query(query string) (bool, error) {
	// Try formula
	url := fmt.Sprintf("https://formulae.brew.sh/api/formula/%s.json", query)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	// Try cask
	url = fmt.Sprintf("https://formulae.brew.sh/api/cask/%s.json", query)
	resp, err = http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK
	return ok, nil
}
