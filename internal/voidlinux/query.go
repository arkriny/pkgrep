package voidlinux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseBody struct {
	Data []json.RawMessage `json:"data"`
}

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://xq-api.voidlinux.org/v1/query/x86_64?q=%s", query)
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

	ok := len(r.Data) != 0
	return ok, nil
}
