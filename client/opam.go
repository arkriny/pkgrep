package client

import (
	"fmt"
	"net/http"

	"github.com/arkriny/pkgrep/httputil"
)

type Opam struct {
	HTTPClient *http.Client
}

func (c Opam) Name() string {
	return "opam"
}

func (c Opam) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://opam.ocaml.org/packages/%s/", query)
	return httputil.GetCheckOK(c.HTTPClient, url)
}
