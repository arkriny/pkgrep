package client

import (
	"fmt"
	"net/http"
)

type Opam struct {
	HTTPClient *http.Client
}

func (c Opam) Name() string {
	return "opam"
}

func (c Opam) Query(query string) (bool, error) {
	url := fmt.Sprintf("https://opam.ocaml.org/packages/%s/", query)
	return getCheckOK(c.HTTPClient, url)
}
