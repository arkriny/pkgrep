package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Chocolatey struct {
	HTTPClient *http.Client
}

func (c Chocolatey) Name() string {
	return "Chocolatey"
}

func (c Chocolatey) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://community.chocolatey.org/packages/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
