package nixpkgs

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Nixpkgs"
}

func (c Client) Query(query string) (bool, error) {
	if len(query) < 2 {
		return false, nil
	}
	url := fmt.Sprintf("https://api.github.com/repos/NixOS/nixpkgs/contents/pkgs/by-name/%s/%s", query[:2], query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
