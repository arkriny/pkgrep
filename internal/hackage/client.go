package hackage

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Hackage"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://hackage.haskell.org/package/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
