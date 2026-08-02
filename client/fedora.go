package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Fedora struct {
	HTTPClient *http.Client
}

func (c Fedora) Name() string {
	return "Fedora"
}

func (c Fedora) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://packages.fedoraproject.org/pkgs/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
