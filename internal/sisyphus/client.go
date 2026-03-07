package sisyphus

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct{}

func (Client) Query(query string) (bool, error) {
	contains, err := httputil.GetBodyContains(fmt.Sprintf("https://packages.altlinux.org/en/sisyphus/srpms/%s", query), "404: That page no exists")
	return !contains, err
}
