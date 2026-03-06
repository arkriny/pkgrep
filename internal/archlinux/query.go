package archlinux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseBody struct {
	Results []json.RawMessage `json:"results"`
}

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://archlinux.org/packages/search/json/?name=%s", query)
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

	ok := len(r.Results) != 0
	return ok, nil
}
