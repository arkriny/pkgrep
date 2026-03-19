package voidlinux

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Void"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/void-linux/void-packages/contents/srcpkgs/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
