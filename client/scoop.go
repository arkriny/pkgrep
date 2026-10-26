package client

import (
	"fmt"
	"net/http"
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
	ok, err := getCheckOK(c.HTTPClient, url)
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	// Try Extras repository
	url = fmt.Sprintf("https://api.github.com/repos/ScoopInstaller/Extras/contents/bucket/%s.json", query)
	return getCheckOK(c.HTTPClient, url)
}
