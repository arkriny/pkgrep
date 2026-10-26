package client

import (
	"fmt"
	"net/http"
)

type Aosc struct {
	HTTPClient *http.Client
}

func (c Aosc) Name() string {
	return "AOSC"
}

func (c Aosc) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://packages.aosc.io/packages/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
