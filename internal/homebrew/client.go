package homebrew

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Homebrew"
}

func (c Client) Query(query string) (bool, error) {
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
