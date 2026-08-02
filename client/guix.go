package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Guix struct {
	HTTPClient *http.Client
}

func (c Guix) Name() string {
	return "Guix"
}

func (c Guix) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://hpc.guix.info/package/%s", query)
	contains, err := httputil.GetBodyContains(c.HTTPClient, url, "<title>Guix-HPC — Oops!</title>")
	return !contains, err
}
