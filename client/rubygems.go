package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Rubygems struct {
	HTTPClient *http.Client
}

func (c Rubygems) Name() string {
	return "RubyGems"
}

func (c Rubygems) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://rubygems.org/api/v1/gems/%s.json", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
