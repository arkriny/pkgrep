package cratesio

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "crates.io"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://crates.io/api/v1/crates/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
