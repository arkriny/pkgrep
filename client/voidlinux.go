package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Voidlinux struct {
	HTTPClient *http.Client
}

func (c Voidlinux) Name() string {
	return "Void"
}

func (c Voidlinux) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/void-linux/void-packages/contents/srcpkgs/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
