package opensuse

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct{}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://software.opensuse.org/package/%s", query)
	contains, err := httputil.GetBodyContains(url, "not found...")
	return !contains, err
}
