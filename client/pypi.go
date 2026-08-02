package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Pypi struct {
	HTTPClient *http.Client
}

func (c Pypi) Name() string {
	return "PyPI"
}

func (c Pypi) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
