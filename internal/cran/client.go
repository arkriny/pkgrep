package cran

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "CRAN"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://crandb.r-pkg.org/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
