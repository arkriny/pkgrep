package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Nuget struct {
	HTTPClient *http.Client
}

func (c Nuget) Name() string {
	return "NuGet"
}

func (c Nuget) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://www.nuget.org/packages/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
