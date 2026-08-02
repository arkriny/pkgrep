package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Clojars struct {
	HTTPClient *http.Client
}

func (c Clojars) Name() string {
	return "Clojars"
}

func (c Clojars) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://clojars.org/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
