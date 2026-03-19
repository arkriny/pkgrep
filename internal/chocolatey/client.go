package chocolatey

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/internal/httputil"
)

type Client struct {
	HTTPClient *http.Client
}

func (c Client) Name() string {
	return "Chocolatey"
}

func (c Client) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://community.chocolatey.org/packages/%s", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
