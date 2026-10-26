package client

import (
	"fmt"
	"net/http"
)

type Voidlinux struct {
	HTTPClient *http.Client
}

func (c Voidlinux) Name() string {
	return "Void"
}

func (c Voidlinux) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/void-linux/void-packages/contents/srcpkgs/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
