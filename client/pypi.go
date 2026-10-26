package client

import (
	"fmt"
	"net/http"
)

type Pypi struct {
	HTTPClient *http.Client
}

func (c Pypi) Name() string {
	return "PyPI"
}

func (c Pypi) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", query)
	return getCheckOK(c.HTTPClient, url)
}
