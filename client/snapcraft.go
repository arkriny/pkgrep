package client

import (
	"fmt"
	"net/http"
)

type Snapcraft struct {
	HTTPClient *http.Client
}

func (c Snapcraft) Name() string {
	return "Snapcraft"
}

func (c Snapcraft) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://snapcraft.io/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
