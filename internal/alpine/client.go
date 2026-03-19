package alpine

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Alpine"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkgs.alpinelinux.org/packages?branch=edge&name=%s", query)
	contains, err := httputil.GetBodyContains(c.HTTPClient, url, "No matching packages found...")
	return !contains, err
}
