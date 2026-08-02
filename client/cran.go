package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Cran struct {
	HTTPClient *http.Client
}

func (c Cran) Name() string {
	return "CRAN"
}

func (c Cran) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://crandb.r-pkg.org/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
