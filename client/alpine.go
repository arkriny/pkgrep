package client

import (
	"fmt"
	"net/http"
)

type Alpine struct {
	HTTPClient *http.Client
}

func (c Alpine) Name() string {
	return "Alpine"
}

func (c Alpine) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkgs.alpinelinux.org/packages?branch=edge&name=%s", query)
	contains, err := getBodyContains(c.HTTPClient, url, "No matching packages found...")
	return !contains, err
}
