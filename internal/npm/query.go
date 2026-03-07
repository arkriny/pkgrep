package npm

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://registry.npmjs.com/%s", query)
	return httputil.GetCheckOK(url)
}
