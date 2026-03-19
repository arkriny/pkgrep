package melpa

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "MELPA"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://melpa.org/packages/%s-badge.svg", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
