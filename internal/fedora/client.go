package fedora

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Fedora"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://packages.fedoraproject.org/pkgs/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
