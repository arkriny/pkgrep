package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Scoop struct {
	HTTPClient *http.Client
}

func (c Scoop) Name() string {
	return "Scoop"
}

func (c Scoop) Query(query string) (bool, error) {
	// Try Main repository
	url := fmt.Sprintf("https://api.github.com/repos/ScoopInstaller/Main/contents/bucket/%s.json", query)
	ok, err := httputil.GetCheckOK(c.HTTPClient, url)
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	// Try Extras repository
	url = fmt.Sprintf("https://api.github.com/repos/ScoopInstaller/Extras/contents/bucket/%s.json", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
