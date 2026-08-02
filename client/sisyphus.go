package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Sisyphus struct {
	HTTPClient *http.Client
}

func (c Sisyphus) Name() string {
	return "Sisyphus"
}

func (c Sisyphus) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://packages.altlinux.org/en/sisyphus/srpms/%s", query)
	contains, err := httputil.GetBodyContains(c.HTTPClient, url, "404: That page no exists")
	return !contains, err
}
