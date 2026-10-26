package client

import (
	"fmt"
	"net/http"
)

type Melpa struct {
	HTTPClient *http.Client
}

func (c Melpa) Name() string {
	return "MELPA"
}

func (c Melpa) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://melpa.org/packages/%s-badge.svg", query)
	return getCheckOK(c.HTTPClient, url)
}
