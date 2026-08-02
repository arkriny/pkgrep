package main

import (
	"encoding/hex"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/arkriny/pkgrep/client"
)

type testcase struct {
	querier Querier
	pkg     string
}

var testUserAgent = "pkgrep-test/0 (https://github.com/arkriny/pkgrep; mailto:arkriny@gmail.com)"

var testHTTPClient = &http.Client{
	Transport: &UserAgentRoundTripper{
		RoundTripper: http.DefaultTransport,
		UserAgent:    testUserAgent,
	},
	Timeout: time.Minute,
}

var tests = []testcase{
	{client.Alpine{testHTTPClient}, "alpine-base"},
	{client.Aosc{testHTTPClient}, "kernel-base"},
	{client.Archlinux{testHTTPClient}, "linux"},
	{client.Aur{testHTTPClient}, "google-chrome"},
	{client.Chocolatey{testHTTPClient}, "go"},
	{client.Clojars{testHTTPClient}, "core.typed"},
	{client.Cran{testHTTPClient}, "ggplot2"},
	{client.Cratesio{testHTTPClient}, "syn"},
	{client.Debian{testHTTPClient}, "linux-base"},
	{client.Dub{testHTTPClient}, "vibe-d"},
	{client.Fedora{testHTTPClient}, "linux-firmware"},
	{client.Guix{testHTTPClient}, "go"},
	{client.Hackage{testHTTPClient}, "ghc"},
	{client.Hex{testHTTPClient}, "phoenix"},
	{client.Homebrew{testHTTPClient}, "go"},
	{client.Julia{testHTTPClient}, "plots"},
	{client.Kali{testHTTPClient}, "linux"},
	{client.Luarocks{testHTTPClient}, "lua-cjson"},
	{client.Macports{testHTTPClient}, "go"},
	{client.Melpa{testHTTPClient}, "magit"},
	{client.Nixpkgs{testHTTPClient}, "home-manager"},
	{client.Npm{testHTTPClient}, "npm"},
	{client.Nuget{testHTTPClient}, "Azure.Core"},
	{client.Opam{testHTTPClient}, "ocaml"},
	{client.Pkggodev{testHTTPClient}, "http"},
	{client.Pubdev{testHTTPClient}, "http"},
	{client.Pypi{testHTTPClient}, "pip"},
	{client.Rubygems{testHTTPClient}, "rails"},
	{client.Scoop{testHTTPClient}, "go"},
	{client.Sisyphus{testHTTPClient}, "firmware-linux"},
	{client.Snapcraft{testHTTPClient}, "go"},
	{client.Ubuntu{testHTTPClient}, "linux-firmware"},
	{client.Voidlinux{testHTTPClient}, "linux-firmware"},
}

func Test(t *testing.T) {
	t.Parallel()

	b := make([]byte, 10)
	rand.Read(b)
	randomPackage := hex.EncodeToString(b)

	for _, tt := range tests {
		name := tt.querier.Name()
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			{
				found, err := tt.querier.Query(tt.pkg)
				if err != nil {
					t.Fatalf("%s: %v", name, err)
				}

				if !found {
					t.Errorf("%s: %s not found", name, tt.pkg)
				}
			}
			{
				found, err := tt.querier.Query(randomPackage)
				if err != nil {
					t.Fatalf("%s: %v", name, err)
				}

				if found {
					t.Errorf("%s: %s found", name, randomPackage)
				}
			}
		})
	}
}
