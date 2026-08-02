package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Macports struct {
	HTTPClient *http.Client
}

func (c Macports) Name() string {
	return "MacPorts"
}

func (c Macports) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://ports.macports.org/api/v1/ports/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
