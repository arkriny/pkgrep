package macports

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "MacPorts"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://ports.macports.org/api/v1/ports/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
