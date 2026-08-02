package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Julia struct {
	HTTPClient *http.Client
}

func (c Julia) Name() string {
	return "Julia"
}

func (c Julia) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://juliapackages.com/p/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
