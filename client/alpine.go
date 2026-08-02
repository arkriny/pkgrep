package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Alpine struct {
	HTTPClient *http.Client
}

func (c Alpine) Name() string {
	return "Alpine"
}

func (c Alpine) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkgs.alpinelinux.org/packages?branch=edge&name=%s", query)
	contains, err := httputil.GetBodyContains(c.HTTPClient, url, "No matching packages found...")
	return !contains, err
}
