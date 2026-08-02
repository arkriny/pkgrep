package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Aur struct {
	HTTPClient *http.Client
}

func (c Aur) Name() string {
	return "AUR"
}

func (c Aur) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://aur.archlinux.org/rpc/v5/info/%s", query)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	var r struct {
		ResultCount int `json:"resultcount"`
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false, err
	}

	ok := r.ResultCount != 0
	return ok, nil
}
