package ubuntu

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct{}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.launchpad.net/devel/ubuntu/+source/%s", query)
	return httputil.GetCheckOK(url)
}
