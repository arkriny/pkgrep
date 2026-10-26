package client

import (
	"fmt"
	"net/http"
)

type Pubdev struct {
	HTTPClient *http.Client
}

func (c Pubdev) Name() string {
	return "pub.dev"
}

func (c Pubdev) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pub.dev/api/packages/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
