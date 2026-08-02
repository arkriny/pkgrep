package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Melpa struct {
	HTTPClient *http.Client
}

func (c Melpa) Name() string {
	return "MELPA"
}

func (c Melpa) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://melpa.org/packages/%s-badge.svg", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
