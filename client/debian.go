package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Debian struct {
	HTTPClient *http.Client
}

func (c Debian) Name() string {
	return "Debian"
}

func (c Debian) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://sources.debian.org/api/src/%s", query)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	var r struct {
		Error *json.RawMessage `json:"error"`
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false, err
	}

	ok := r.Error == nil
	return ok, nil
}
