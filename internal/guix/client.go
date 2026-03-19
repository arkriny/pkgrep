package guix

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Guix"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://hpc.guix.info/package/%s", query)
	contains, err := httputil.GetBodyContains(c.HTTPClient, url, "<title>Guix-HPC — Oops!</title>")
	return !contains, err
}
