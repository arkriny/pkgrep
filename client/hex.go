package client

import (
	"fmt"
	"net/http"
)

type Hex struct {
	HTTPClient *http.Client
}

func (c Hex) Name() string {
	return "Hex"
}

func (c Hex) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://hex.pm/api/packages/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
