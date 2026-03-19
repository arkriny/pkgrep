package julia

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Julia"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://juliapackages.com/p/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
