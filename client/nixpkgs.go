package client

import (
	"fmt"
	"net/http"
)

type Nixpkgs struct {
	HTTPClient *http.Client
}

func (c Nixpkgs) Name() string {
	return "Nixpkgs"
}

func (c Nixpkgs) Query(query string) (bool, error) {
	if len(query) < 2 {
		return false, nil
	}
	url := fmt.Sprintf("https://api.github.com/repos/NixOS/nixpkgs/contents/pkgs/by-name/%s/%s", query[:2], query)
	return getCheckOK(c.HTTPClient, url)
}
