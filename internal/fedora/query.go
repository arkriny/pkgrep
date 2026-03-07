package fedora

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://packages.fedoraproject.org/pkgs/%s", query)
	return httputil.GetCheckOK(url)
}
