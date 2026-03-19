package dub

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "DUB"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://code.dlang.org/packages/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
