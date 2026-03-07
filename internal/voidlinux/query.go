package voidlinux

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/void-linux/void-packages/contents/srcpkgs/%s", query)
	return httputil.GetCheckOK(url)
}
