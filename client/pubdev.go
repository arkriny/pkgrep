package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Pubdev struct {
	HTTPClient *http.Client
}

func (c Pubdev) Name() string {
	return "pub.dev"
}

func (c Pubdev) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pub.dev/api/packages/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
