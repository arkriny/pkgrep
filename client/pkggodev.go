package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Pkggodev struct {
	HTTPClient *http.Client
}

func (c Pkggodev) Name() string {
	return "pkg.go.dev"
}

func (c Pkggodev) Query(query string) (bool, error) {
	// filter to prevent matches within paths and descriptions.
	filter := fmt.Sprintf(`hasSuffix(modulePath, "/%s")`, query)
	url := fmt.Sprintf("https://pkg.go.dev/v1/search?q=%s&filter=%s&limit=1", query, url.QueryEscape(filter))
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	var r struct {
		Total int `json:"total"`
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false, err
	}

	ok := r.Total != 0
	return ok, nil
}
