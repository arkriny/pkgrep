package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Dub struct {
	HTTPClient *http.Client
}

func (c Dub) Name() string {
	return "DUB"
}

func (c Dub) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://code.dlang.org/packages/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
