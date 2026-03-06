package nixpkgs

import (
	"fmt"
	"net/http"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/NixOS/nixpkgs/contents/pkgs/by-name/%s/%s", query[:2], query)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK
	return ok, nil
}
