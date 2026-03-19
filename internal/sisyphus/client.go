package sisyphus

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Sisyphus"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://packages.altlinux.org/en/sisyphus/srpms/%s", query)
	contains, err := httputil.GetBodyContains(c.HTTPClient, url, "404: That page no exists")
	return !contains, err
}
