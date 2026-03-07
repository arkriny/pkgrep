package snapcraft

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://snapcraft.io/%s", query)
	return httputil.GetCheckOK(url)
}
