package client

import (
	"fmt"
	"net/http"
)

type Kali struct {
	HTTPClient *http.Client
}

func (c Kali) Name() string {
	return "Kali"
}

func (c Kali) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkg.kali.org/pkg/%s", query)
	return getCheckOK(c.HTTPClient, url)
}
