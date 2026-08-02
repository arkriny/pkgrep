package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Archlinux struct {
	HTTPClient *http.Client
}

func (c Archlinux) Name() string {
	return "Arch"
}

type responseBody struct {
	Results []json.RawMessage `json:"results"`
}

func (c Archlinux) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://archlinux.org/packages/search/json/?name=%s", query)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	var r responseBody
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false, err
	}

	ok := len(r.Results) != 0
	return ok, nil
}
