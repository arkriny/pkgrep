package macports

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct{}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://ports.macports.org/api/v1/ports/%s", query)
	return httputil.GetCheckOK(url)
}
