package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Homebrew struct {
	HTTPClient *http.Client
}

func (c Homebrew) Name() string {
	return "Homebrew"
}

func (c Homebrew) Query(query string) (bool, error) {
	// Try formula
	url := fmt.Sprintf("https://formulae.brew.sh/api/formula/%s.json", query)
	ok, err := httputil.GetCheckOK(c.HTTPClient, url)
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	// Try cask
	url = fmt.Sprintf("https://formulae.brew.sh/api/cask/%s.json", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
