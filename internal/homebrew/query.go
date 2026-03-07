package homebrew

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

func Query(query string) (bool, error) {
	// Try formula
	url := fmt.Sprintf("https://formulae.brew.sh/api/formula/%s.json", query)
	ok, err := httputil.GetCheckOK(url)
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	// Try cask
	url = fmt.Sprintf("https://formulae.brew.sh/api/cask/%s.json", query)
	return httputil.GetCheckOK(url)
}
