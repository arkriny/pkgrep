package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Kali struct {
	HTTPClient *http.Client
}

func (c Kali) Name() string {
	return "Kali"
}

func (c Kali) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkg.kali.org/pkg/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
