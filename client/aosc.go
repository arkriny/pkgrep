package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Aosc struct {
	HTTPClient *http.Client
}

func (c Aosc) Name() string {
	return "AOSC"
}

func (c Aosc) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://packages.aosc.io/packages/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
