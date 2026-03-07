package debian

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct{}

type responseBody struct {
	Error *json.RawMessage `json:"error"`
}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://sources.debian.org/api/src/%s", query)
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

	ok := r.Error == nil
	return ok, nil
}
