package client

import (
	"fmt"
	"net/http"
)

type Clojars struct {
	HTTPClient *http.Client
}

func (c Clojars) Name() string {
	return "Clojars"
}

func (c Clojars) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://clojars.org/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
