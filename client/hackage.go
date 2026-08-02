package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Hackage struct {
	HTTPClient *http.Client
}

func (c Hackage) Name() string {
	return "Hackage"
}

func (c Hackage) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://hackage.haskell.org/package/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
