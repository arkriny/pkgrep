package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Ubuntu struct {
	HTTPClient *http.Client
}

func (c Ubuntu) Name() string {
	return "Ubuntu"
}

func (c Ubuntu) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.launchpad.net/devel/ubuntu/+source/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
