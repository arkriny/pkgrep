package client

import (
	"fmt"
	"net/http"
)

type Npm struct {
	HTTPClient *http.Client
}

func (c Npm) Name() string {
	return "NPM"
}

func (c Npm) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://registry.npmjs.com/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
