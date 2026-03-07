package nixpkgs

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct{}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/NixOS/nixpkgs/contents/pkgs/by-name/%s/%s", query[:2], query)
	return httputil.GetCheckOK(url)
}
