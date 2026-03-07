package rubygems

import (
	"fmt"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct{}

func (Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://rubygems.org/api/v1/gems/%s.json", query)
	return httputil.GetCheckOK(url)
}
