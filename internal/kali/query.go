package kali

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

func Query(query string) (bool, error) {
	url := fmt.Sprintf("https://pkg.kali.org/pkg/%s", query)
	return httputil.GetCheckOK(url)
}
