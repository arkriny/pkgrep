package clojars

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Clojars"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://clojars.org/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
