package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Snapcraft struct {
	HTTPClient *http.Client
}

func (c Snapcraft) Name() string {
	return "Snapcraft"
}

func (c Snapcraft) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://snapcraft.io/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
