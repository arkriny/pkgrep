package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Cratesio struct {
	HTTPClient *http.Client
}

func (c Cratesio) Name() string {
	return "crates.io"
}

func (c Cratesio) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://crates.io/api/v1/crates/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
