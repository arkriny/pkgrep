package client

import (
	"fmt"
	"net/http"
)

type Julia struct {
	HTTPClient *http.Client
}

func (c Julia) Name() string {
	return "Julia"
}

func (c Julia) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://juliapackages.com/p/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
