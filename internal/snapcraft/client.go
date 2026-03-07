package snapcraft

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct{}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://snapcraft.io/%s", query)
	return httputil.GetCheckOK(url)
}
