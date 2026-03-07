package aur

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct{}

type responseBody struct {
	ResultCount int `json:"resultcount"`
}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://aur.archlinux.org/rpc/v5/info/%s", query)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var r responseBody
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false, err
	}

	ok := r.ResultCount != 0
	return ok, nil
}
